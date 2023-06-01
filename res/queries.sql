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
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4326)) as location
FROM nlwkn_water_rights.e_usage_locations;

-- name: get-water-rights-by-reality
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4326)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE real = $1 OR real IS NULL;

-- name: get-water-rights-by-state
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4326)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE active = $1 OR active IS NULL;

-- name: get-water-rights-by-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4236)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE ST_CONTAINS(ST_Collect(ARRAY((SELECT geom FROM geodata.shapes WHERE key = any($2)))), st_transform(location, 4236));

-- name: get-water-rights-by-reality-and-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4236)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE (real = $1 OR real IS NULL)
AND ST_CONTAINS(ST_Collect(ARRAY((SELECT geom FROM geodata.shapes WHERE key = any($2)))), st_transform(location, 4236));

-- name: get-water-rights-by-state-and-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4236)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE (active = $1 OR active IS NULL)
AND ST_CONTAINS(ST_Collect(ARRAY((SELECT geom FROM geodata.shapes WHERE key = any($2)))), st_transform(location, 4236));

-- name: get-water-rights-by-reality-and-state
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4236)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE (real = $1 OR real IS NULL)
AND active = $2 OR active IS NULL;

-- name: get-water-rights-by-reality-state-and-location
SELECT id, water_right, active, real, name, ST_ASGEOJSON(st_transform(location, 4236)) as location
FROM nlwkn_water_rights.e_usage_locations
WHERE( real = $1 OR real IS NULL)
AND (active = $2 OR active IS NULL)
AND ST_CONTAINS(ST_Collect(ARRAY((SELECT geom FROM geodata.shapes WHERE key = any($2)))), st_transform(location, 4236));

-- name: get-water-right-details
SELECT
    id, no, ext_id, file_ref, legal_title, state, subject, address, annotation,
    bailee, date_of_change, valid, granting_authority, registering_authority, water_authority
FROM nlwkn_water_rights.water_rights
WHERE no = $1::int
LIMIT 1;

-- name: get-detailed-locations
SELECT
    id,
    water_right,
    name,
    no,
    active,
    ST_ASGEOJSON(st_transform(location, 4236)) as location,
    (CASE
        WHEN basin_no is NULL THEN null
        WHEN basin_no is not NULL THEN
            jsonb_build_object('key', (basin_no).key, 'name', (basin_no).name)
        END) as basin_no,
    county,
    (CASE
        WHEN eu_survey_area is NULL THEN null
        WHEN eu_survey_area is not NULL THEN
            jsonb_build_object('key', (eu_survey_area).key, 'name', (eu_survey_area).name)
        END) as eu_survey_area,
    field,
    groundwater_volume,
    legal_scope,
    local_sub_district,
    (CASE
        WHEN maintenance_association is NULL THEN null
        WHEN maintenance_association is not NULL THEN
            jsonb_build_object('key', (maintenance_association).key, 'name', (maintenance_association).name)
        END) as maintenance_association,
    (CASE
        WHEN municipal_area is NULL THEN null
        WHEN municipal_area is not NULL THEN
            jsonb_build_object('key', (municipal_area).key, 'name', (municipal_area).name)
        END) as municipal_area,
    plot,
    real,
    rivershed,
    serial_no,
    (CASE
        WHEN top_map_1_25000 is NULL THEN null
        WHEN top_map_1_25000 is not NULL THEN
            jsonb_build_object('key', (top_map_1_25000).key, 'name', (top_map_1_25000).name)
        END) as top_map_1_25000,
    water_body,
    flood_area,
    water_protection_area,
    array_to_json(withdrawal_rate) as withdrawal_rate,
    array_to_json(fluid_discharge) as fluid_discharge,
    (CASE
        WHEN irrigation_area is NULL THEN null
        WHEN irrigation_area is not NULL THEN
            jsonb_build_object('amount', (irrigation_area).amount, 'unit', (irrigation_area).unit)
        END) as irrigation_area ,
    array_to_json(rain_supplement) as rain_supplement
FROM nlwkn_water_rights.e_usage_locations
WHERE water_right = $1::int