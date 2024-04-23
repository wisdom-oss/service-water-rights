-- name: water-rights
SELECT *
FROM water_rights.rights;

-- name: get-current-wr
SELECT internal_id
FROM water_rights.current_rights
WHERE water_right_number = $1;

-- name: get-water-right
SELECT *
FROM water_rights.rights
WHERE id = $1;

-- name: get-locations
SELECT *
FROM water_rights.usage_locations;

-- name: get-water-right-usage-locations
SELECT *
FROM water_rights.usage_locations
WHERE water_right = $1;