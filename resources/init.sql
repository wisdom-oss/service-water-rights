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
            key  numeric,
            name text,
            per  interval
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
        CREATE TYPE water_rights.land_record AS
        (
            district text,
            field    int8
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END;
$$;

-- name: 06
DO
$$
    BEGIN
        CREATE TYPE water_rights.value_range AS
        (
            district text,
            field    int8
        );
    EXCEPTION
        WHEN DUPLICATE_OBJECT THEN NULL;
    END;
$$;

-- name: 07
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

-- name: 08
CREATE TABLE IF NOT EXISTS water_rights.rights
(
    id                    int8 NOT NULL PRIMARY KEY,
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
    granting_authority    text                            DEFAULT NULL
);

-- name: 09
CREATE TABLE IF NOT EXISTS water_rights.usage_locations
(
    id                      bigserial PRIMARY KEY,
    no                      int8                                                  DEFAULT NULL,
    serial                  text                                                  DEFAULT NULL,
    water_right             int8 REFERENCES water_rights.rights (id) MATCH SIMPLE DEFAULT NULL,
    active                  bool                                                  DEFAULT NULL,
    real                    bool                                                  DEFAULT NULL,
    name                    text                                                  DEFAULT NULL,
    legal_purpose           text[2]                                               DEFAULT NULL,
    map_excerpt             water_rights.numeric_keyed_value                      DEFAULT NULL,
    municipal_area          water_rights.numeric_keyed_value                      DEFAULT NULL,
    county                  text                                                  DEFAULT NULL,
    land_record             water_rights.land_record                              DEFAULT NULL,
    plot                    text                                                  DEFAULT NULL,
    maintenance_association water_rights.numeric_keyed_value                      DEFAULT NULL,
    eu_survey_area          water_rights.numeric_keyed_value                      DEFAULT NULL,
    catchment_area_code     water_rights.numeric_keyed_value                      DEFAULT NULL,
    regulation_citation     text                                                  DEFAULT NULL,
    withdrawal_rates        water_rights.rate[]                                   DEFAULT NULL,
    pumping_rates           water_rights.rate[]                                   DEFAULT NULL,
    injection_rates         water_rights.rate[]                                   DEFAULT NULL,
    waste_water_flow_volume water_rights.rate[]                                   DEFAULT NULL,
    river_basin             text                                                  DEFAULT NULL,
    groundwater_body        text                                                  DEFAULT NULL,
    water_body              text                                                  DEFAULT NULL,
    flood_area              text                                                  DEFAULT NULL,
    water_protection_area   text                                                  DEFAULT NULL,
    dam_target_levels       water_rights.dam_target                               DEFAULT NULL,
    fluid_discharge         water_rights.rate[]                                   DEFAULT NULL,
    rain_supplement         water_rights.rate[]                                   DEFAULT NULL,
    irrigation_area         water_rights.quantity                                 DEFAULT NULL,
    ph_values               numrange                                              DEFAULT NULL,
    injection_limits        text[][2]                                             DEFAULT NULL,
    location                geometry                                              DEFAULT NULL
);