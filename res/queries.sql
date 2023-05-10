-- name: create-schema
create schema if not exists nlwkn_water_rights;

-- name: create-type-interval-rate
create type nlwkn_water_rights.interval_rate as
(
    amount   integer,
    unit     varchar,
    duration interval
);

-- name: create-type-numeric-keyed-name
create type nlwkn_water_rights.numeric_keyed_name as
(
    key  integer,
    name varchar
);

-- name: create-type-rate
create type nlwkn_water_rights.rate as
(
    amount integer,
    unit   varchar
);

-- name: create-type-water-right-state
create type nlwkn_water_rights.water_right_state as enum ('aktiv', 'inaktiv', 'Wasserbuchblatt');


-- name: create-water-rights-table
create table water_rights
(
    id                    serial,
    no                    integer not null,
    ext_id                varchar,
    file_ref              varchar,
    legal_title           varchar,
    state                 nlwkn_water_rights.water_right_state,
    subject               varchar,
    address               varchar,
    annotation            varchar,
    bailee                varchar,
    date_of_change        timestamp,
    valid                 daterange,
    granting_authority    varchar,
    registering_authority varchar,
    water_authority       varchar
);