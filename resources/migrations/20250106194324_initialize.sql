-- +goose Up
-- +goose StatementBegin
-- create schema
CREATE SCHEMA IF NOT EXISTS water_rights;

-- create data types
DO $$ BEGIN
CREATE TYPE water_rights.quantity AS ("value" NUMERIC, unit TEXT);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
CREATE TYPE water_rights.rate AS ("value" NUMERIC, unit TEXT, per INTERVAL);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
CREATE TYPE water_rights."legal_department" AS ENUM('A', 'B', 'C', 'D', 'E', 'F', 'K', 'L');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
CREATE TYPE water_rights.numeric_keyed_value AS ("key" NUMERIC, "name" TEXT);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
CREATE TYPE water_rights.land_record AS ("district" TEXT, "field" int8, "fallback" TEXT);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
CREATE TYPE water_rights.injection_limit AS ("substance" TEXT, "quantity" water_rights.quantity);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$ BEGIN
CREATE TYPE water_rights.dam_target AS (
    "default" water_rights.quantity,
    "steady" water_rights.quantity,
    "max" water_rights.quantity
);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- create tables
CREATE TABLE IF NOT EXISTS
    water_rights.rights (
        id bigserial NOT NULL,
        water_right_number int8 NOT NULL,
        external_identifier TEXT NULL,
        file_reference TEXT NULL,
        legal_departments water_rights."legal_department" NULL,
        holder TEXT NULL,
        address TEXT NULL,
        subject TEXT NULL,
        legal_title TEXT NULL,
        status TEXT NULL,
        valid_from date NULL,
        valid_until date NULL,
        initially_granted date NULL,
        last_change date NULL,
        water_authority TEXT NULL,
        registering_authority TEXT NULL,
        granting_authority TEXT NULL,
        annotation TEXT NULL,
        CONSTRAINT rights_pkey PRIMARY KEY (id)
    );

CREATE TABLE IF NOT EXISTS
    water_rights.current_rights (
        water_right_number int8 NOT NULL,
        internal_id int8 NULL,
        deleted timestamptz NULL,
        CONSTRAINT current_rights_pkey PRIMARY KEY (water_right_number)
    );

CREATE TABLE IF NOT EXISTS
    water_rights.usage_locations (
        id bigserial NOT NULL,
        "no" int8 NULL,
        serial TEXT NULL,
        water_right int8 NOT NULL,
        "legal_department" water_rights."legal_department" NOT NULL,
        active bool NULL,
        "real" bool NULL,
        "name" TEXT NULL,
        legal_purpose _text NULL,
        map_excerpt water_rights.numeric_keyed_value NULL,
        municipal_area water_rights.numeric_keyed_value NULL,
        county TEXT NULL,
        land_record water_rights.land_record NULL,
        plot TEXT NULL,
        maintenance_association water_rights.numeric_keyed_value NULL,
        eu_survey_area water_rights.numeric_keyed_value NULL,
        catchment_area_code water_rights.numeric_keyed_value NULL,
        regulation_citation TEXT NULL,
        withdrawal_rates water_rights."rate" NULL,
        pumping_rates water_rights."rate" NULL,
        injection_rates water_rights."rate" NULL,
        waste_water_flow_volume water_rights."rate" NULL,
        river_basin TEXT NULL,
        groundwater_body TEXT NULL,
        water_body TEXT NULL,
        flood_area TEXT NULL,
        water_protection_area TEXT NULL,
        dam_target_levels water_rights.dam_target NULL,
        fluid_discharge water_rights."rate" NULL,
        rain_supplement water_rights."rate" NULL,
        irrigation_area water_rights.quantity NULL,
        ph_values numrange NULL,
        injection_limits water_rights."injection_limit" NULL,
        "location" public.geometry (POINT, 25832) DEFAULT NULL::geometry NULL,
        CONSTRAINT usage_locations_pkey PRIMARY KEY (id)
    );

-- +goose StatementEnd