-- name: create-water-right-schema
CREATE TABLE IF NOT EXISTS water_rights.water_rights
(
    id                    bigint NOT NULL PRIMARY KEY,
    rights_holder         text                       DEFAULT NULL,
    valid_from            date                       DEFAULT '-infinity'::date,
    valid_until           date                       DEFAULT 'infinity'::date,
    status                text                       DEFAULT NULL,
    legal_title           text                       DEFAULT NULL,
    water_authority       text                       DEFAULT NULL,
    registering_authority text                       DEFAULT NULL,
    granting_authority    text                       DEFAULT NULL,
    first_granted         date                       DEFAULT NULL,
    date_of_change        date                       DEFAULT NULL,
    file_reference        text                       DEFAULT NULL,
    external_identifier   text                       DEFAULT NULL,
    subject               text                       DEFAULT NULL,
    address               text                       DEFAULT NULL,
    legal_departments     water_rights.departments[] DEFAULT NULL
);


-- name: usage-locations
SELECT id,
    water_right,
    name,
    no,
    active,
    real,
    ST_AsGeoJSON(ST_TRANSFORM(location, 4326))::jsonb AS location
-- basin_no, county,
-- eu_survey_area, field,
-- groundwater_volume,
-- legal_scope, local_sub_district, maintenance_association, municipal_area, plot, real, rivershed, serial_no,
-- top_map_1_25000, water_body,
-- flood_area, water_protection_area, withdrawal_rate, fluid_discharge, irrigation_area, rain_supplement
FROM nlwkn_water_rights.e_usage_locations;

-- name: water-rights
SELECT id,
    no,
    ext_id,
    file_ref,
    legal_title,
    state,
    subject,
    address,
    annotation,
    bailee,
    date_of_change,
    valid,
    granting_authority,
    registering_authority,
    water_authority
FROM nlwkn_water_rights.water_rights;


-- name: filter-locations
ST_CONTAINS
(ST_COLLECT(ARRAY ((SELECT geom FROM geodata.shapes WHERE KEY = ANY ($1)))),ST_TRANSFORM(LOCATION, 4326));

-- name: filter-reality
    (REAL = $1 OR REAL IS NULL);

-- name: filter-state
    (active = $1 OR active IS NULL);