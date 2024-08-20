-- name: water-rights
SELECT *
FROM water_rights.rights;


-- name: get-current-wr
SELECT internal_id
FROM water_rights.current_rights
WHERE
    water_right_number = $1;

-- name: get-water-right
SELECT *
FROM water_rights.rights
WHERE
    id = $1;

-- name: get-locations
SELECT *, st_asewkb(st_transform(location, 4326)) AS location_ewkb
FROM water_rights.usage_locations;


-- name: get-water-right-usage-locations
SELECT *, st_asewkb(st_transform(location, 4326)) AS location_ewkb
FROM water_rights.usage_locations
WHERE
    water_right = $1;

-- name: get-withdrawal-rates
SELECT withdrawal_rates
FROM water_rights.usage_locations
WHERE id IN (SELECT internal_id FROM water_rights.current_rights) AND st_within(st_transform(location, 4326), st_transform(st_geomfromgeojson($1::text), 4326::integer)) AND withdrawal_rates is not NULL;

-- name: filter-locations
ST_CONTAINS(ST_COLLECT(ARRAY ((SELECT geom FROM geodata.shapes WHERE KEY = ANY ($1)))),ST_TRANSFORM(LOCATION, 4326));

-- name: filter-reality
(real = $1 OR real IS NULL);

-- name: filter-state
(active = $1 OR active IS NULL);