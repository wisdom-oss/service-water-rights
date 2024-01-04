-- name: filter-location
ST_CONTAINS(ST_COLLECT(ARRAY((SELECT geom FROM geodata.shapes WHERE key = any($1)))),ST_TRANSFORM(location, 4236));

-- name: filter-reality
real = $1 OR real IS NULL;

-- name: filter-state
active = $1 OR active IS NULL;

-- name: extended-usage-locations
SELECT
    id, water_right, name, no, active, ST_ASGEOJSON(ST_TRANSFORM(location, 4236)) as location, basin_no, county,
    eu_survey_area, field,
    groundwater_volume,
    legal_scope, local_sub_district, maintenance_association, municipal_area, plot, real, rivershed, serial_no,
    top_map_1_25000, water_body,
    flood_area, water_protection_area, withdrawal_rate, fluid_discharge, irrigation_area, rain_supplement
FROM
    nlwkn_water_rights.e_usage_locations;

-- name: usage-locations
SELECT
    id, water_right, name, no, active, ST_ASGEOJSON(ST_TRANSFORM(location, 4236)) as location, real
FROM
    nlwkn_water_rights.e_usage_locations;

-- name: water-rights
SELECT
    id, no, ext_id, file_ref, legal_title, state, subject, address, annotation,
    bailee, date_of_change, valid, granting_authority, registering_authority, water_authority
FROM nlwkn_water_rights.water_rights;