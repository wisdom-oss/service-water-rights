import asyncio
import logging
import os
import socket
import sys

import amqp_rpc_client
import orjson
import pydantic
import requests

import configuration
import enums
import models.amqp
import models.internal
import tools

bind = f"0.0.0.0:{configuration.ServiceConfiguration().http_port}"
workers = 1
limit_request_line = 0
limit_request_fields = 0
limit_request_field_size = 0
worker_class = "uvicorn.workers.UvicornWorker"
max_requests = 0
timeout = 0


def on_starting(server):
    _service_configuration = configuration.ServiceConfiguration()
    logging.basicConfig(
        format="[%(asctime)s] [%(process)d] [%(levelname)s] %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S %z",
        level=_service_configuration.logging_level,
        force=True,
    )
    # %% Validate the AMQP configuration and message broker reachability
    try:
        _amqp_configuration = configuration.AMQPConfiguration()
    except pydantic.ValidationError:
        logging.critical(
            "Unable to read the service registry related settings. Please refer to "
            "the documentation for further instructions: "
            "AMQP_CONFIGURATION_INVALID"
        )
        sys.exit(1)
    _amqp_configuration.dsn.port = 5672 if _amqp_configuration.dsn.port is None else int(_amqp_configuration.dsn.port)
    _message_broker_reachable = asyncio.run(
        tools.is_host_available(_amqp_configuration.dsn.host, _amqp_configuration.dsn.port)
    )
    if not _message_broker_reachable:
        logging.error(
            f"The message broker is currently not reachable on {_amqp_configuration.dsn.host}:"
            f"{_amqp_configuration.dsn.port}"
        )
        sys.exit(2)
    # %% Check if the configured service scope is available
    # Create an amqp client
    _amqp_client = amqp_rpc_client.Client(amqp_dsn=_amqp_configuration.dsn, mute_pika=False)
    service_scope = models.internal.ServiceScope.parse_file("./configuration/scope.json")
    # Query if the scope already exists
    _scope_check_request = models.amqp.CheckScopeRequest(value=service_scope.value)
    _scope_check_request_id = _amqp_client.send(
        _scope_check_request.json(), _amqp_configuration.authorization_exchange, "authorization-service"
    )
    _scope_check_response_bytes = _amqp_client.await_response(_scope_check_request_id)
    _scope_check_response: dict = orjson.loads(_scope_check_response_bytes)
    logging.info("Got following response from the AMQP Auth Service: %s", _scope_check_response)
    # Check if the scope check response contains any of the known error keys
    if set(_scope_check_response.keys()).issubset({"httpCode", "httpError", "error", "errorName", "errorDescription"}):
        # Since the scope check response contains an error request the scope to be created
        _scope_create_request = models.amqp.CreateScopeRequest(
            name=service_scope.name, description=service_scope.description, value=service_scope.value
        )
        _scope_create_request_id = _amqp_client.send(
            _scope_create_request.json(), _amqp_configuration.authorization_exchange
        )
        _scope_create_response_bytes = _amqp_client.await_response(_scope_create_request_id)
        _scope_create_response: dict = orjson.loads(_scope_create_response_bytes)
        logging.info("Got following response from the AMQP Auth Service: %s", _scope_create_response)
        if set(_scope_create_response.keys()).issubset(
            {"httpCode", "httpError", "error", "errorName", "errorDescription"}
        ):
            logging.critical(
                "Unable to create the scope which shall be used by the service:\n%s", _scope_create_response
            )
            sys.exit(3)
        logging.info("Successfully created the scope that shall be used by this service")
    # Set the value for the used security scope internally
    os.environ["CONFIG_SECURITY_SCOPE"] = service_scope.value
    # %% Validate the database settings and reachability
    try:
        _database_configuration = configuration.DatabaseConfiguration()
    except pydantic.ValidationError:
        logging.critical(
            "Unable to read the service registry related settings. Please refer to "
            "the documentation for further instructions: "
            "DATABASE_CONFIGURATION_INVALID"
        )
        sys.exit(1)
    logging.info("Checking the connection to the database")
    _database_configuration.dsn.port = (
        5432 if _database_configuration.dsn.port is None else int(_database_configuration.dsn.port)
    )
    _database_available = asyncio.run(
        tools.is_host_available(
            host=_database_configuration.dsn.host, port=_database_configuration.dsn.port, timeout=10
        )
    )
    if not _database_available:
        logging.critical(
            "The database is not available. Since this service requires an access to the database the service will "
            "not start"
        )
        sys.exit(2)
    try:
        _gateway_information = configuration.KongGatewayInformation()
    except pydantic.ValidationError:
        logging.critical(
            "Unable to read the information about the Kong API Gateway. Please refer to the documentation for further "
            "instructions: KONG_INFORMATION_INVALID "
        )
        sys.exit(1)
    _gateway_reachable = asyncio.run(
        tools.is_host_available(_gateway_information.hostname, _gateway_information.admin_port)
    )
    if not _gateway_reachable:
        logging.critical("The api gateway is not available. Since the service needs to register itself on the ")
        sys.exit(2)


def when_ready(server):
    # %% Register at the Kong gateway
    _gateway_information = configuration.KongGatewayInformation()
    _service_settings = configuration.ServiceConfiguration()
    logging.debug("Read the following information about the gateway:\n%s", _gateway_information.json(indent=2))
    # Try to get information about the upstream
    upstream_information_request = tools.query_kong(
        f"/upstreams/upstream_{_service_settings.name}", enums.HTTPMethod.GET
    )
    if upstream_information_request.status_code == 404:
        logging.warning("No upstream for this service found. Creating a new upstream in the API gateway...")
        new_upstream_information = {"name": f"upstream_{_service_settings.name}"}
        upstream_creation = tools.query_kong(
            f"/upstreams/",
            enums.HTTPMethod.POST,
            data=new_upstream_information,
        )
        if upstream_creation.status_code == 201:
            logging.info(
                "Created a new upstream for this service:\n%s",
                orjson.dumps(upstream_creation.json(), option=orjson.OPT_INDENT_2),
            )
        else:
            logging.debug(
                f"Received a {upstream_creation.status_code} from the gateway:\n%s",
                orjson.dumps(upstream_creation.json(), option=orjson.OPT_INDENT_2),
            )
    elif upstream_information_request.status_code == 200:
        logging.debug(
            "Found the following upstream information for this service:\n%s",
            orjson.dumps(upstream_information_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
        )

    service_information_request = tools.query_kong(f"/services/service_{_service_settings.name}", enums.HTTPMethod.GET)
    if service_information_request.status_code == 404:
        logging.warning("No service entry found for this service. Creating a new entry in the API gateway...")
        new_service_information = {
            "name": f"service_{_service_settings.name}",
            "host": f"upstream_{_service_settings.name}",
        }
        service_creation_request = tools.query_kong(f"/services/", enums.HTTPMethod.POST, new_service_information)
        if service_creation_request.status_code == 201:
            logging.info(
                "Created a new entry for this service:\n%s",
                orjson.dumps(service_creation_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
            )
    elif service_information_request.status_code == 200:
        logging.debug(
            "Found the following information for this service:\n%s",
            orjson.dumps(service_information_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
        )
    route_information_request = tools.query_kong(
        f"/services/service_" f"{_service_settings.name}/routes/{_gateway_information.service_path_slug}",
        method=enums.HTTPMethod.GET,
    )
    if route_information_request.status_code == 404:
        logging.warning("No route is configured for this service. Creating a new route definition for this service...")
        route_creation_request_data = {
            "paths[]": f"/{_gateway_information.service_path_slug}",
            "name": _gateway_information.service_path_slug,
        }
        route_creation_request = tools.query_kong(
            f"/services/service_{_service_settings.name}/routes/",
            method=enums.HTTPMethod.POST,
            data=route_creation_request_data,
        )
        if route_creation_request.status_code == 201:
            logging.info(
                "Created a new route for this service:\n%s",
                orjson.dumps(route_creation_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
            )
    # Determine the ip address of the service container
    ip_address = socket.gethostbyname(socket.gethostname())
    # Request information about the available targets
    upstream_target_information_request = tools.query_kong(
        f"/upstreams/upstream_{_service_settings.name}/targets", enums.HTTPMethod.GET
    )
    upstream_target_information = upstream_target_information_request.json()
    logging.debug(
        "Got following upstream information:\n%s",
        orjson.dumps(upstream_target_information, option=orjson.OPT_INDENT_2).decode("utf-8"),
    )
    container_listed = any(
        [
            target["target"] == f"{ip_address}:{_service_settings.http_port}"
            for target in upstream_target_information["data"]
        ]
    )
    if not container_listed:
        upstream_target_creation_data = {"target": f"{ip_address}:{_service_settings.http_port}"}
        upstream_creation_request = tools.query_kong(
            f"/upstreams/upstream_{_service_settings.name}/targets",
            enums.HTTPMethod.POST,
            upstream_target_creation_data,
        )
        if upstream_creation_request.status_code == 201:
            logging.info(
                "Created a new upstream target for this service:\n%s",
                orjson.dumps(upstream_creation_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
            )
    consumer_information_request = tools.query_kong("/consumers", enums.HTTPMethod.GET)
    consumer_exists = any(
        [consumer["custom_id"] == "authorization-service" for consumer in consumer_information_request.json()["data"]]
    )
    _consumer_id = None
    if not consumer_exists:
        consumer_creation_request_data = {"custom_id": "authorization-service"}
        consumer_creation_request = tools.query_kong(
            "/consumers", data=consumer_creation_request_data, method=enums.HTTPMethod.POST
        )
        if consumer_creation_request.status_code == 201:
            logging.info(
                "Created new consumer for this service:\n%s",
                orjson.dumps(consumer_creation_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
            )
            _consumer_id = consumer_creation_request.json()["id"]
    else:
        _consumer_id = [
            consumer["id"]
            for consumer in consumer_information_request.json()["data"]
            if consumer["custom_id"] == "authorization-service"
        ][0]
    consumer_credential_information_request = tools.query_kong(
        f"/consumers/{_consumer_id}/oauth2", method=enums.HTTPMethod.GET
    )
    consumer_credentials_exists = any(
        [
            credential["consumer"]["id"] == _consumer_id
            for credential in consumer_credential_information_request.json()["data"]
        ]
    )
    if not consumer_credentials_exists:
        consumer_credential_creation_request_data = {
            "name": "Authorization Module",
            "redirect_uris": "http://localhost/authenticated",
        }
        consumer_credential_creation_request = tools.query_kong(
            f"/consumers/{_consumer_id}/oauth2",
            data=consumer_credential_creation_request_data,
            method=enums.HTTPMethod.POST,
        )
        if consumer_credential_creation_request.status_code == 201:
            logging.info(
                "Created new consumer credentials for this service:\n%s",
                orjson.dumps(consumer_credential_creation_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
            )
            credential_file = open("/.credential_id", "wt")
            credential_file.write(consumer_credential_creation_request.json()["id"])
    else:
        _credential_id = [
            credential["id"]
            for credential in consumer_credential_information_request.json()["data"]
            if credential["consumer"]["id"] == _consumer_id
        ][0]
        credential_file = open("/.credential_id", "wt")
        credential_file.write(_credential_id)
    plugin_information_request = tools.query_kong(
        f"/services/service_{_service_settings.name}/plugins", method=enums.HTTPMethod.GET
    )
    plugins = [plugin for plugin in plugin_information_request.json()["data"]]
    oauth2_configured = any(["oauth2" == plugin["name"] for plugin in plugins])
    if not oauth2_configured:
        plugin_creation_data = {
            "name": "oauth2",
            "config.enable_authorization_code": "false",
            "config.enable_client_credentials": "true",
            "config.enable_implicit_grant": "false",
            "config.enable_password_grant": "false",
            "config.hide_credentials": "true",
            "config.accept_http_if_already_terminated": "true",
            "config.global_credentials": "true",
        }
        plugin_creation_request = tools.query_kong(
            f"/services/service_{_service_settings.name}/plugins",
            data=plugin_creation_data,
            method=enums.HTTPMethod.POST,
        )
        if plugin_creation_request.status_code == 201:
            logging.info(
                "Created new OAuth2 plugin for this service:\n%s",
                orjson.dumps(plugin_creation_request.json(), option=orjson.OPT_INDENT_2).decode("utf-8"),
            )


def on_exit(server):
    _service_settings = configuration.ServiceConfiguration()
    _gateway_information = configuration.KongGatewayInformation()
    ip_address = socket.gethostbyname(socket.gethostname())
    upstream_deletion_request = requests.delete(
        f"http://{_gateway_information.hostname}:{_gateway_information.admin_port}/upstreams/upstream_"
        f"{_service_settings.name}/targets/{ip_address}:{_service_settings.http_port}",
    )
