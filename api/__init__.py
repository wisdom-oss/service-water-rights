"""Package containing the code which will be the API later on"""
import datetime
import email.utils
import hashlib
import http
import typing

import amqp_rpc_client
import fastapi
import geoalchemy2.functions
import orjson
import pytz as pytz
import redis
import sqlalchemy.dialects
import sqlalchemy.exc
import starlette.middleware.gzip
import starlette.responses

import api.handler
import configuration
import database
import database.tables
import exceptions
import models.internal
import tools
from api import security

# %% API Setup
service = fastapi.FastAPI()
service.add_exception_handler(exceptions.APIException, api.handler.handle_api_error)
service.add_exception_handler(fastapi.exceptions.RequestValidationError, api.handler.handle_request_validation_error)
service.add_exception_handler(sqlalchemy.exc.IntegrityError, api.handler.handle_integrity_error)

# %% Configurations
_security_configuration = configuration.SecurityConfiguration()
_service_configuration = configuration.ServiceConfiguration()
_redis_configuration = configuration.RedisConfiguration()

if _security_configuration.scope_string_value is None:
    service_scope = models.internal.ServiceScope.parse_file("./configuration/scope.json")
    _security_configuration.scope_string_value = service_scope.value

# %% Global Clients
_redis_client = redis.Redis.from_url(_redis_configuration.dsn)


# %% Middlewares
@service.middleware("http")
async def etag_comparison(request: fastapi.Request, call_next):
    """
    A middleware which will hash the request path and all parameters transferred to this
    microservice and will check if the hash matches the one of the ETag which was sent to the
    microservice. Furthermore, it will take the generated hash and append it to the response to
    allow caching

    :param request: The incoming request
    :type request: fastapi.Request
    :param call_next: The next call after this middleware
    :type call_next: callable
    :return: The result of the next call after this middle ware
    :rtype: fastapi.Response
    """
    # Access all parameters used for creating the hash
    path = request.url.path
    query_parameter = dict(request.query_params)
    # Now iterate through all query parameters and make sure they are sorted if they are lists
    for key, value in dict(query_parameter).items():
        # Now check if the value is a list
        if isinstance(value, list):
            query_parameter[key] = sorted(value)
    query_dict = {
        "request_path": path,
        "request_query_parameter": query_parameter,
    }
    query_data = orjson.dumps(query_dict, option=orjson.OPT_SORT_KEYS)
    # Now create a hashsum of the query data
    query_hash = hashlib.sha3_256(query_data).hexdigest()
    # Create redis keys for later usage
    response_cache_key = _service_configuration.name + ".data." + query_hash
    response_change_cache_key = _service_configuration.name + ".last_change." + query_hash
    # Now access the headers of the request and check for the If-None-Match Header
    e_tag = request.headers.get("If-None-Match", None)
    last_known_update = request.headers.get("If-Modified-Since", _redis_client.get(response_change_cache_key))
    if last_known_update is None:
        last_known_update = datetime.datetime.fromtimestamp(0, tz=pytz.UTC)
    else:
        if type(last_known_update) is bytes:
            last_known_update = email.utils.parsedate_to_datetime(last_known_update.decode("utf-8"))
        else:
            last_known_update = email.utils.parsedate_to_datetime(last_known_update)
    # Get the last update of the schema from which the service gets its data from
    last_database_modification = tools.get_last_schema_update("nlwkn_water_rights", database.engine)
    data_changed = last_known_update < last_database_modification
    if data_changed:
        response: starlette.responses.StreamingResponse = await call_next(request)
        if response.status_code == 200:
            _redis_client.set(response_change_cache_key, email.utils.format_datetime(last_database_modification))
            response_content = [chunk async for chunk in response.body_iterator][0].decode()
            _redis_client.set(response_cache_key, response_content)
            response.headers.append("ETag", f"{query_hash}")
            response.headers.append("Last-Modified", email.utils.format_datetime(last_database_modification))

            return fastapi.Response(
                content=response_content,
                headers={"E-Tag": query_hash, "Last-Modified": email.utils.format_datetime(last_database_modification)},
                media_type="text/json",
            )
        return response
    if _redis_client.get(response_cache_key) is None:
        response: starlette.responses.StreamingResponse = await call_next(request)
        if response.status_code == 200:
            _redis_client.set(response_change_cache_key, email.utils.format_datetime(last_database_modification))
            response_content = [chunk async for chunk in response.body_iterator][0]
            _redis_client.set(response_cache_key, response_content)
            response.headers.append("ETag", f"{query_hash}")
            response.headers.append("Last-Modified", email.utils.format_datetime(last_database_modification))
            return fastapi.Response(
                content=response_content,
                headers={"E-Tag": query_hash, "Last-Modified": email.utils.format_datetime(last_database_modification)},
                media_type="text/json",
            )
        return response
    else:
        return fastapi.Response(
            content=_redis_client.get(response_cache_key),
            headers={"E-Tag": query_hash, "Last-Modified": email.utils.format_datetime(last_database_modification)},
            media_type="text/json",
        )


# %% Routes
@service.get("/")
async def get(
    in_area: None | list[str] = fastapi.Query(default=None, alias="in"),
    is_active: None | bool = fastapi.Query(default=None, alias="is_active"),
    is_real: None | bool = fastapi.Query(default=None, alias="is_real"),
    user: models.internal.UserAccount
    | bool = fastapi.Security(security.is_authorized_user, scopes=[_security_configuration.scope_string_value]),
):
    # Build a tuple with the available parameter
    available_parameter = (in_area is not None, is_active is not None, is_real is not None)
    # the columns which are queried
    query_columns = [
        database.tables.e_usage_locations.c.id,
        database.tables.e_usage_locations.c.water_right,
        database.tables.e_usage_locations.c.active,
        database.tables.e_usage_locations.c.real,
        database.tables.e_usage_locations.c.name,
        sqlalchemy.cast(
            geoalchemy2.functions.ST_AsGeoJSON(
                geoalchemy2.functions.ST_Transform(database.tables.e_usage_locations.c.location, 4326)
            ),
            sqlalchemy.dialects.postgresql.JSONB,
        ).label("geojson"),
    ]
    match available_parameter:
        case (False, False, False):
            location_filter = None
        case (False, False, True):
            location_filter = sqlalchemy.or_(
                database.tables.e_usage_locations.c.real == is_real, database.tables.e_usage_locations.c.real == None
            )
        case (False, True, False):
            location_filter = sqlalchemy.or_(
                database.tables.e_usage_locations.c.active == is_active,
                database.tables.e_usage_locations.c.active == None,
            )
        case (True, False, False):
            location_filter = sqlalchemy.or_(
                *[
                    geoalchemy2.functions.ST_Contains(
                        sqlalchemy.select(
                            [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                            database.tables.shapes.c.key == k,
                        ),
                        database.tables.e_usage_locations.c.location,
                    )
                    for k in in_area
                ],
            )
        case (False, True, True):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.e_usage_locations.c.real == is_real,
                    database.tables.e_usage_locations.c.real == None,
                ),
                sqlalchemy.or_(
                    database.tables.e_usage_locations.c.active == is_active,
                    database.tables.e_usage_locations.c.active == None,
                ),
            )
        case (True, False, True):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.e_usage_locations.c.real == is_real,
                    database.tables.e_usage_locations.c.real == None,
                ),
                sqlalchemy.or_(
                    *[
                        geoalchemy2.functions.ST_Contains(
                            sqlalchemy.select(
                                [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                                database.tables.shapes.c.key == k,
                            ),
                            database.tables.e_usage_locations.c.location,
                        )
                        for k in in_area
                    ],
                ),
            )
        case (True, True, False):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.e_usage_locations.c.active == is_active,
                    database.tables.e_usage_locations.c.active == None,
                ),
                sqlalchemy.or_(
                    *[
                        geoalchemy2.functions.ST_Contains(
                            sqlalchemy.select(
                                [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                                database.tables.shapes.c.key == k,
                            ),
                            database.tables.e_usage_locations.c.location,
                        )
                        for k in in_area
                    ],
                ),
            )
        case (True, True, True):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.e_usage_locations.c.real == is_real,
                    database.tables.e_usage_locations.c.real == None,
                ),
                sqlalchemy.or_(
                    database.tables.e_usage_locations.c.active == is_active,
                    database.tables.e_usage_locations.c.active == None,
                ),
                sqlalchemy.or_(
                    *[
                        geoalchemy2.functions.ST_Contains(
                            sqlalchemy.select(
                                [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                                database.tables.shapes.c.key == k,
                            ),
                            database.tables.e_usage_locations.c.location,
                        )
                        for k in in_area
                    ],
                ),
            )
        case _:
            location_filter = None
    location_query = sqlalchemy.select(query_columns, location_filter)
    locations = database.engine.execute(location_query).all()
    if len(locations) == 0:
        return fastapi.Response(status_code=http.HTTPStatus.NO_CONTENT)
    return locations


@service.get("/details/{water_right_number}")
async def get_details(
    water_right_number: int | None = fastapi.Path(default=...),
    user: models.internal.UserAccount
    | bool = fastapi.Security(security.is_authorized_user, scopes=[_security_configuration.scope_string_value]),
):
    water_right_query_columns = [
        database.tables.water_rights.c.id,
        database.tables.water_rights.c.no,
        database.tables.water_rights.c.ext_id.label("externalId"),
        database.tables.water_rights.c.file_ref.label("fileReference"),
        database.tables.water_rights.c.legal_title.label("legalTitle"),
        database.tables.water_rights.c.state,
        database.tables.water_rights.c.subject,
        database.tables.water_rights.c.address,
        database.tables.water_rights.c.annotation,
        database.tables.water_rights.c.bailee,
        database.tables.water_rights.c.date_of_change.label("dateOfChange"),
        database.tables.water_rights.c.valid,
        database.tables.water_rights.c.granting_authority.label("grantingAuthority"),
        database.tables.water_rights.c.registering_authority.label("registeringAuthority"),
        database.tables.water_rights.c.water_authority.label("waterAuthority"),
    ]

    e_usage_location_query_columns = [
        database.tables.e_usage_locations.c.id,
        database.tables.e_usage_locations.c.water_right,
        database.tables.e_usage_locations.c.name,
        database.tables.e_usage_locations.c.no,
        database.tables.e_usage_locations.c.active,
        sqlalchemy.cast(
            geoalchemy2.functions.ST_AsGeoJSON(
                geoalchemy2.functions.ST_Transform(database.tables.e_usage_locations.c.location, 4326)
            ),
            sqlalchemy.dialects.postgresql.JSONB,
        ).label("location"),
        database.tables.e_usage_locations.c.basin_no,
        database.tables.e_usage_locations.c.county,
        database.tables.e_usage_locations.c.eu_survey_area,
        database.tables.e_usage_locations.c.field,
        database.tables.e_usage_locations.c.groundwater_volume,
        database.tables.e_usage_locations.c.legal_scope,
        database.tables.e_usage_locations.c.local_sub_district,
        database.tables.e_usage_locations.c.maintenance_association,
        database.tables.e_usage_locations.c.municipal_area,
        database.tables.e_usage_locations.c.plot,
        database.tables.e_usage_locations.c.real,
        database.tables.e_usage_locations.c.rivershed,
        database.tables.e_usage_locations.c.serial_no,
        database.tables.e_usage_locations.c.top_map_1_25000,
        database.tables.e_usage_locations.c.water_body,
        database.tables.e_usage_locations.c.flood_area,
        database.tables.e_usage_locations.c.water_protection_area,
        database.tables.e_usage_locations.c.withdrawal_rate,
        database.tables.e_usage_locations.c.fluid_discharge,
        database.tables.e_usage_locations.c.irrigation_area,
        database.tables.e_usage_locations.c.rain_supplement,
    ]

    water_right_query = sqlalchemy.select(
        water_right_query_columns, database.tables.water_rights.c.no == water_right_number
    )
    water_right = database.engine.execute(water_right_query).first()
    # Pull all e-type usage locations
    e_usage_location_query = f"""SELECT nlwkn_water_rights.e_usage_locations.id,
    nlwkn_water_rights.e_usage_locations.water_right, nlwkn_water_rights.e_usage_locations.name,
    nlwkn_water_rights.e_usage_locations.no, nlwkn_water_rights.e_usage_locations.active, CAST(ST_AsGeoJSON(
    ST_Transform(nlwkn_water_rights.e_usage_locations.location, 4326)) AS JSONB) AS location,
    nlwkn_water_rights.e_usage_locations.basin_no, nlwkn_water_rights.e_usage_locations.county,
     nlwkn_water_rights.e_usage_locations.eu_survey_area, nlwkn_water_rights.e_usage_locations.field,
      nlwkn_water_rights.e_usage_locations.groundwater_volume, nlwkn_water_rights.e_usage_locations.legal_scope,
       nlwkn_water_rights.e_usage_locations.local_sub_district, nlwkn_water_rights.e_usage_locations.maintenance_association,
        nlwkn_water_rights.e_usage_locations.municipal_area, nlwkn_water_rights.e_usage_locations.plot, nlwkn_water_rights.e_usage_locations.real,
         nlwkn_water_rights.e_usage_locations.rivershed, nlwkn_water_rights.e_usage_locations.serial_no,
         nlwkn_water_rights.e_usage_locations.top_map_1_25000,
          nlwkn_water_rights.e_usage_locations.water_body,
         nlwkn_water_rights.e_usage_locations.flood_area, nlwkn_water_rights.e_usage_locations.water_protection_area, nlwkn_water_rights.e_usage_locations.withdrawal_rate,
          nlwkn_water_rights.e_usage_locations.fluid_discharge, nlwkn_water_rights.e_usage_locations.irrigation_area, nlwkn_water_rights.e_usage_locations.rain_supplement
FROM nlwkn_water_rights.e_usage_locations
WHERE nlwkn_water_rights.e_usage_locations.water_right = {water_right_number}"""
    usage_locations = database.engine.execute(e_usage_location_query)
    locations = []
    for usage_location in usage_locations:
        print(usage_location)
        locations.append(usage_location)
    return {**dict(water_right), "locations": locations}
