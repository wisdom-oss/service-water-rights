/*
 This file contains all queries needed to initialize the database for using this
 microservice.
 The queries are safe if the initialization already ran once as the queries
 handle duplicate items by themselves.
 */

-- name: 01
CREATE SCHEMA IF NOT EXISTS water_rights;

-- name: 02
DO
$$
    BEGIN
        CREATE TYPE water_rights.numeric_keyed_value AS
        (
            key  numeric,
            name text
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END
$$;

-- name: 03
DO
$$
    BEGIN
        CREATE TYPE water_rights.quantity AS
        (
            value numeric,
            unit  text
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END
$$;

-- name: 04
DO
$$
    BEGIN
        CREATE TYPE water_rights.rate AS
        (
            value numeric,
            unit  text,
            per   interval
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END
$$;

-- name: 04
DO
$$
    BEGIN
        CREATE TYPE water_rights.dam_target AS
        (
            "default" water_rights.quantity,
            steady    water_rights.quantity,
            max       water_rights.quantity
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END
$$;

-- name: 05
DO
$$
    BEGIN
        -- the attributes "district" and "fields" are mutually exclusive to "fallback"
        CREATE TYPE water_rights.land_record AS
        (
            district text,
            field    int8,
            fallback text
        );
        COMMENT ON TYPE water_rights.land_record IS 'the attributes "district" and "fields" are mutually exclusive to "fallback"';
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END;
$$;

-- name: 07
DO
$$
    BEGIN
        CREATE TYPE water_rights.injection_limit AS
        (
            substance text,
            quantity  water_rights.quantity
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END;
$$;

-- name: 08
DO
$$
    BEGIN
        CREATE TYPE water_rights.legal_department AS ENUM
            (
                'A',
                'B',
                'C',
                'D',
                'E',
                'F',
                'K',
                'L'
                );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END;
$$;

-- name: 09
CREATE TABLE IF NOT EXISTS water_rights.rights
(
    -- id is the internally used id that is automatically created for each water
    -- right
    id                    bigserial PRIMARY KEY,
    -- water_right_number contains the officially issued water right number
    -- by the governing body
    water_right_number    bigint NOT NULL,
    external_identifier   text                            DEFAULT NULL,
    file_reference        text                            DEFAULT NULL,
    legal_departments     water_rights.legal_department[] DEFAULT NULL,
    holder                text                            DEFAULT NULL,
    address               text                            DEFAULT NULL,
    subject               text                            DEFAULT NULL,
    legal_title           text                            DEFAULT NULL,
    status                text                            DEFAULT NULL,
    valid_from            date                            DEFAULT NULL,
    valid_until           date                            DEFAULT NULL,
    initially_granted     date                            DEFAULT NULL,
    last_change           date                            DEFAULT NULL,
    water_authority       text                            DEFAULT NULL,
    registering_authority text                            DEFAULT NULL,
    granting_authority    text                            DEFAULT NULL,
    annotation            text                            DEFAULT NULL
);

-- name: 10
CREATE TABLE IF NOT EXISTS water_rights.usage_locations
(
    id                      bigserial PRIMARY KEY,
    no                      bigint                             DEFAULT NULL,
    serial                  text                             DEFAULT NULL,
    -- water_right points to the internally used id of the water right since the
    -- locations may vary between each water right
    water_right             bigint REFERENCES water_rights.rights (id) MATCH FULL NOT NULL,
    legal_department        water_rights.legal_department                       NOT NULL,
    active                  bool                             DEFAULT NULL,
    real                    bool                             DEFAULT NULL,
    name                    text                             DEFAULT NULL,
    legal_purpose           text[2]                          DEFAULT NULL,
    map_excerpt             water_rights.numeric_keyed_value DEFAULT NULL,
    municipal_area          water_rights.numeric_keyed_value DEFAULT NULL,
    county                  text                             DEFAULT NULL,
    land_record             water_rights.land_record         DEFAULT NULL,
    plot                    text                             DEFAULT NULL,
    maintenance_association water_rights.numeric_keyed_value DEFAULT NULL,
    eu_survey_area          water_rights.numeric_keyed_value DEFAULT NULL,
    catchment_area_code     water_rights.numeric_keyed_value DEFAULT NULL,
    regulation_citation     text                             DEFAULT NULL,
    withdrawal_rates        water_rights.rate[]              DEFAULT NULL,
    pumping_rates           water_rights.rate[]              DEFAULT NULL,
    injection_rates         water_rights.rate[]              DEFAULT NULL,
    waste_water_flow_volume water_rights.rate[]              DEFAULT NULL,
    river_basin             text                             DEFAULT NULL,
    groundwater_body        text                             DEFAULT NULL,
    water_body              text                             DEFAULT NULL,
    flood_area              text                             DEFAULT NULL,
    water_protection_area   text                             DEFAULT NULL,
    dam_target_levels       water_rights.dam_target          DEFAULT NULL,
    fluid_discharge         water_rights.rate[]              DEFAULT NULL,
    rain_supplement         water_rights.rate[]              DEFAULT NULL,
    irrigation_area         water_rights.quantity            DEFAULT NULL,
    ph_values               numrange                         DEFAULT NULL,
    injection_limits        water_rights.injection_limit[]   DEFAULT NULL,
    location                geometry('point', 25832)         DEFAULT NULL
);

-- name: 11
CREATE TABLE IF NOT EXISTS water_rights.current_rights
(
    water_right_number int8 NOT NULL,
    internal_id        int8,
    deleted            timestamptz DEFAULT NULL
);