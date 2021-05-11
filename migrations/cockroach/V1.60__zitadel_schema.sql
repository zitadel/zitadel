CREATE DATABASE zitadel;
use zitadel;

CREATE SCHEMA read_models AUTHORIZATION authz;

CREATE TABLE read_models.orgs (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    org_state SMALLINT,
    sequence BIGINT,

    domain TEXT,
    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE read_models.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);