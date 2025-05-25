-- name: v2_get-water-right
SELECT *
FROM water_rights.rights
WHERE id = $1
    OR water_right_number = $1;

-- name: v2_get-water-right-usage-locations
SELECT *
FROM water_rights.usage_locations
WHERE water_right = $1;