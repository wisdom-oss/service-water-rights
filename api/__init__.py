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
# service.add_exception_handler(fastapi.exceptions.RequestValidationError, api.handler.handle_request_validation_error)
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
    user: str
    | bool = fastapi.Security(security.is_authorized_user, scopes=[_security_configuration.scope_string_value]),
):
    # Build a tuple with the available parameter
    available_parameter = (in_area is not None, is_active is not None, is_real is not None)
    # the columns which are queried
    query_columns = [
        database.tables.E_UsageLocation.id,
        database.tables.E_UsageLocation.water_right,
        database.tables.E_UsageLocation.active,
        database.tables.E_UsageLocation.real,
        database.tables.E_UsageLocation.name,
        sqlalchemy.cast(
            geoalchemy2.functions.ST_AsGeoJSON(
                geoalchemy2.functions.ST_Transform(database.tables.E_UsageLocation.location, 4326)
            ),
            sqlalchemy.dialects.postgresql.JSONB,
        ).label("geojson"),
    ]
    match available_parameter:
        case (False, False, False):
            location_filter = None
        case (False, False, True):
            location_filter = sqlalchemy.or_(
                database.tables.E_UsageLocation.real == is_real, database.tables.E_UsageLocation.real == None
            )
        case (False, True, False):
            location_filter = sqlalchemy.or_(
                database.tables.E_UsageLocation.active == is_active,
                database.tables.E_UsageLocation.active == None,
            )
        case (True, False, False):
            location_filter = sqlalchemy.or_(
                *[
                    geoalchemy2.functions.ST_Contains(
                        sqlalchemy.select(
                            [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                            database.tables.shapes.c.key == k,
                        ),
                        database.tables.E_UsageLocation.location,
                    )
                    for k in in_area
                ],
            )
        case (False, True, True):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.E_UsageLocation.real == is_real,
                    database.tables.E_UsageLocation.real == None,
                ),
                sqlalchemy.or_(
                    database.tables.E_UsageLocation.active == is_active,
                    database.tables.E_UsageLocation.active == None,
                ),
            )
        case (True, False, True):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.E_UsageLocation.real == is_real,
                    database.tables.E_UsageLocation.real == None,
                ),
                sqlalchemy.or_(
                    *[
                        geoalchemy2.functions.ST_Contains(
                            sqlalchemy.select(
                                [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                                database.tables.shapes.c.key == k,
                            ),
                            database.tables.E_UsageLocation.location,
                        )
                        for k in in_area
                    ],
                ),
            )
        case (True, True, False):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.E_UsageLocation.active == is_active,
                    database.tables.E_UsageLocation.active == None,
                ),
                sqlalchemy.or_(
                    *[
                        geoalchemy2.functions.ST_Contains(
                            sqlalchemy.select(
                                [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                                database.tables.shapes.c.key == k,
                            ),
                            database.tables.E_UsageLocation.location,
                        )
                        for k in in_area
                    ],
                ),
            )
        case (True, True, True):
            location_filter = sqlalchemy.and_(
                sqlalchemy.or_(
                    database.tables.E_UsageLocation.real == is_real,
                    database.tables.E_UsageLocation.real == None,
                ),
                sqlalchemy.or_(
                    database.tables.E_UsageLocation.active == is_active,
                    database.tables.E_UsageLocation.active == None,
                ),
                sqlalchemy.or_(
                    *[
                        geoalchemy2.functions.ST_Contains(
                            sqlalchemy.select(
                                [geoalchemy2.functions.ST_Transform(database.tables.shapes.c.geom, 25832)],
                                database.tables.shapes.c.key == k,
                            ),
                            database.tables.E_UsageLocation.location,
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
    user: str
    | bool = fastapi.Security(security.is_authorized_user, scopes=[_security_configuration.scope_string_value]),
):
    database.tables.init()
    water_right_query_columns = [
        database.tables.WaterRight.id,
        database.tables.WaterRight.no,
        database.tables.WaterRight.ext_id.label("externalId"),
        database.tables.WaterRight.file_ref.label("fileReference"),
        database.tables.WaterRight.legal_title.label("legalTitle"),
        database.tables.WaterRight.state,
        database.tables.WaterRight.subject,
        database.tables.WaterRight.address,
        database.tables.WaterRight.annotation,
        database.tables.WaterRight.bailee,
        database.tables.WaterRight.date_of_change.label("dateOfChange"),
        database.tables.WaterRight.valid,
        database.tables.WaterRight.granting_authority.label("grantingAuthority"),
        database.tables.WaterRight.registering_authority.label("registeringAuthority"),
        database.tables.WaterRight.water_authority.label("waterAuthority"),
    ]

    e_usage_location_query_columns = [
        database.tables.E_UsageLocation.id,
        database.tables.E_UsageLocation.water_right.label("waterRight"),
        database.tables.E_UsageLocation.name,
        database.tables.E_UsageLocation.no,
        database.tables.E_UsageLocation.active,
        sqlalchemy.cast(
            geoalchemy2.functions.ST_AsGeoJSON(
                geoalchemy2.functions.ST_Transform(database.tables.E_UsageLocation.location, 4326)
            ),
            sqlalchemy.dialects.postgresql.JSONB,
        ).label("location"),
        sqlalchemy.func.to_json(database.tables.E_UsageLocation.basin_no).label("basinNo"),
        database.tables.E_UsageLocation.county,
        sqlalchemy.func.to_json(database.tables.E_UsageLocation.eu_survey_area).label("euSurveyArea"),
        database.tables.E_UsageLocation.field,
        database.tables.E_UsageLocation.groundwater_volume.label("groundwaterVolume"),
        database.tables.E_UsageLocation.legal_scope.label("legalScope"),
        database.tables.E_UsageLocation.local_sub_district.label("localSubDistrict"),
        sqlalchemy.func.to_json(database.tables.E_UsageLocation.maintenance_association).label(
            "maintenanceAssociation"
        ),
        sqlalchemy.func.to_json(database.tables.E_UsageLocation.municipal_area).label("municipalArea"),
        database.tables.E_UsageLocation.plot,
        database.tables.E_UsageLocation.real,
        database.tables.E_UsageLocation.rivershed,
        database.tables.E_UsageLocation.serial_no.label("serialNo"),
        sqlalchemy.func.to_json(database.tables.E_UsageLocation.top_map_1_25000).label("topMap1to25000"),
        database.tables.E_UsageLocation.water_body.label("waterBody"),
        database.tables.E_UsageLocation.flood_area.label("floodArea"),
        database.tables.E_UsageLocation.water_protection_area.label("waterProtectionArea"),
        sqlalchemy.func.to_json(database.tables.E_UsageLocation.irrigation_area).label("irrigationArea"),
    ]

    e_usage_location_withdrawal_rate_query_text = (
        "SELECT id, amount, unit, extract(epoch from duration)::int as duration from (select id, "
        "(unnest(withdrawal_rate)).amount, (unnest(withdrawal_rate)).unit, "
        "(unnest(withdrawal_rate)).duration from nlwkn_water_rights.e_usage_locations "
        "WHERE id = {0}) sub"
    )
    e_usage_location_fluid_discharge_query_text = (
        "SELECT id, amount, unit, extract(epoch from duration)::int as duration from (select id, "
        "(unnest(fluid_discharge)).amount, (unnest(fluid_discharge)).unit, "
        "(unnest(fluid_discharge)).duration from nlwkn_water_rights.e_usage_locations "
        "WHERE id = {0}) sub"
    )
    e_usage_location_rain_supplement_query_text = (
        "SELECT id, amount, unit, extract(epoch from duration)::int as duration from (select id, "
        "(unnest(rain_supplement)).amount, (unnest(rain_supplement)).unit, "
        "(unnest(rain_supplement)).duration from nlwkn_water_rights.e_usage_locations "
        "WHERE id = {0}) sub"
    )
    water_right_query = sqlalchemy.select(*water_right_query_columns).where(
        database.tables.WaterRight.no == water_right_number
    )
    water_right = database.engine.execute(water_right_query).first()
    if water_right is None:
        raise exceptions.APIException(
            error_code="NO_SUCH_WATER_RIGHT",
            error_title="No water right known by this number",
            error_description="There is no water right in the database that matches this number",
            http_status=http.HTTPStatus.NOT_FOUND,
        )
    water_right_dict = dict(water_right)
    water_right_dict.update({"valid": {"from": water_right["valid"].lower, "until": water_right["valid"].upper}})
    # Pull all e-type usage locations
    e_usage_location_query = sqlalchemy.select(
        e_usage_location_query_columns, database.tables.E_UsageLocation.water_right == water_right_number
    )
    usage_locations = database.engine.execute(e_usage_location_query).all()
    _locations = []
    for usage_location in usage_locations:
        _location = dict(usage_location)
        _withdrawal_rates = []
        _fluid_discharges = []
        _rain_supplements = []
        withdrawal_rates = database.engine.execute(
            sqlalchemy.text(e_usage_location_withdrawal_rate_query_text.format(_location["id"]))
        ).all()
        fluid_discharges = database.engine.execute(
            sqlalchemy.text(e_usage_location_fluid_discharge_query_text.format(_location["id"]))
        ).all()
        rain_supplements = database.engine.execute(
            sqlalchemy.text(e_usage_location_rain_supplement_query_text.format(_location["id"]))
        ).all()
        for rate in withdrawal_rates:
            _withdrawal_rates.append(
                {
                    "amount": rate["amount"],
                    "unit": rate["unit"].strip().replace("'", ""),
                    "duration": rate["duration"],
                }
            )
        for rate in fluid_discharges:
            _fluid_discharges.append(
                {
                    "amount": rate["amount"],
                    "unit": rate["unit"].strip().replace("'", ""),
                    "duration": rate["duration"],
                }
            )
        for rate in rain_supplements:
            _rain_supplements.append(
                {
                    "amount": rate["amount"],
                    "unit": rate["unit"].strip().replace("'", ""),
                    "duration": rate["duration"],
                }
            )
        _location.update(
            {
                "withdrawalRates": _withdrawal_rates,
                "fluidDischarge": _fluid_discharges,
                "rainSupplement": _rain_supplements,
            }
        )
        _locations.append(_location)
    return {**water_right_dict, "locations": _locations}
