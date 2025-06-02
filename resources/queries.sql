-- name: water-rights
SELECT
    *
FROM
    water_rights.rights;

-- name: get-current-wr
SELECT
    internal_id
FROM
    water_rights.current_rights
WHERE
    water_right_number = $1
    OR internal_id = $1;

-- name: get-water-right
SELECT
    *
FROM
    water_rights.rights
WHERE
    id = $1;

-- name: get-locations
SELECT
    *
FROM
    water_rights.usage_locations;

-- name: get-water-right-usage-locations
SELECT
    *
FROM
    water_rights.usage_locations
WHERE
    water_right = $1;

-- name: get-withdrawal-rates
SELECT
    withdrawal_rates
FROM
    water_rights.usage_locations
WHERE
    active = true
    AND id IN (
        SELECT
            internal_id
        FROM
            water_rights.current_rights
    )
    AND ST_Within (ST_Transform (location, 4326), $1)
    AND withdrawal_rates is not NULL;