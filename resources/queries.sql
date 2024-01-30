-- name: water-rights
SELECT *
FROM water_rights.rights;

/*
 * Filters for the water rights
 */
-- name: filter-water-right-legal-department
(legal_departments && ANY ($1));
-- name: filter-water-right-ids
(id = any($1));

-- name: usage-locations
SELECT *, st_asgeojson(st_transform(location, 4326))::jsonb AS location
FROM water_rights.usage_locations;

/*
 * Filters for the usage locations
 */
-- name: filter-usage-location-known-geometry
st_contains(
    st_collect(
        ARRAY (
            SELECT geom FROM geodata.shapes WHERE KEY = ANY ($1)
        )
    ),
    st_transform(location, 4326)
);
-- name: filter-usage-location-ids
(id = any($1));
-- name: filter-usage-location-water-rights
(water_right = any($1));
-- name: filter-usage-location-reality
(real = $1 OR real IS NULL);
-- name: filter-usage-location-active-state
(active = $1 OR active IS NULL);
-- name: filter-usage-location-custom-wkt-geometry
st_contains(
    st_setsrid(
        st_geomfromewkt($1), 4326
    ),
    st_transform(location, 4326)
);
-- name: filter-usage-location-custom-geojson-geometry
st_contains(
    st_geomfromgeojson($1),
    st_transform(location, 4326)
)
