-- name: create-schema
create schema if not exists nlwkn_water_rights;

-- name: create-type-interval-rate
create type nlwkn_water_rights.interval_rate as
(
    amount   integer,
    unit     varchar,
    duration interval
);

-- name: create-type-numeric-keyed-name
create type nlwkn_water_rights.numeric_keyed_name as
(
    key  integer,
    name varchar
);

-- name: create-type-rate
create type nlwkn_water_rights.rate as
(
    amount integer,
    unit   varchar
);

-- name: create-type-water-right-state
create type nlwkn_water_rights.water_right_state as enum ('aktiv', 'inaktiv', 'Wasserbuchblatt');


-- name: create-water-rights-table
create table water_rights
(
    id                    serial,
    no                    integer not null,
    ext_id                varchar,
    file_ref              varchar,
    legal_title           varchar,
    state                 nlwkn_water_rights.water_right_state,
    subject               varchar,
    address               varchar,
    annotation            varchar,
    bailee                varchar,
    date_of_change        timestamp,
    valid                 daterange,
    granting_authority    varchar,
    registering_authority varchar,
    water_authority       varchar
);

-- name: create-table-usage-locations
create table nlwkn_water_rights.usage_locations
(
    id                      serial,
    water_right             integer not null,
    name                    varchar,
    no                      integer,
    active                  boolean,
    location                geometry(Point, 25832),
    basin_no                nlwkn_water_rights.numeric_keyed_name,
    county                  varchar,
    eu_survey_area          nlwkn_water_rights.numeric_keyed_name,
    field                   integer,
    groundwater_volume      varchar,
    legal_scope             varchar,
    local_sub_district      varchar,
    maintenance_association nlwkn_water_rights.numeric_keyed_name,
    municipal_area          nlwkn_water_rights.numeric_keyed_name,
    plot                    varchar,
    real                    boolean,
    rivershed               varchar,
    serial_no               varchar,
    top_map_1_25000         nlwkn_water_rights.numeric_keyed_name,
    water_body              varchar
);

-- name: create-derived-table
create table nlwkn_water_rights.e_usage_locations
(
    flood_area            varchar,
    water_protection_area varchar,
    withdrawal_rate       nlwkn_water_rights.interval_rate[],
    fluid_discharge       nlwkn_water_rights.interval_rate[],
    irrigation_area       nlwkn_water_rights.rate,
    rain_supplement       nlwkn_water_rights.interval_rate[]
)
    inherits (nlwkn_water_rights.usage_locations);


-- name: get-all-water-rights
SELECT id, water_right, active, real, name
FROM nlwkn_water_rights.e_usage_locations;

-- name: get-water-rights-by-reality
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE real = $1 OR real IS NULL;

-- name: get-water-rights-by-state
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE active = $1 OR active IS NULL;

-- name: get-water-rights-by-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE ST_CONTAINS((SELECT geom FROM geodata.shapes WHERE key = any($1)), location);

-- name: get-water-rights-by-reality-and-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE real = $1 OR real IS NULL
AND ST_CONTAINS((SELECT geom FROM geodata.shapes WHERE key = any($2)), location);

-- name: get-water-rights-by-state-and-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE active = $1 OR active IS NULL
AND ST_CONTAINS((SELECT geom FROM geodata.shapes WHERE key = any($2)), location);

-- name: get-water-rights-by-reality-and-state
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE real = $1 OR real IS NULL
AND active = $2 OR active IS NULL;

-- name: get-water-rights-by-reality-state-and-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(location)
FROM nlwkn_water_rights.e_usage_locations
WHERE real = $1 OR real IS NULL
AND active = $2 OR active IS NULL
AND ST_CONTAINS((SELECT geom FROM geodata.shapes WHERE key = any($3)), location);