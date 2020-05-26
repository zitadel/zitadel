BEGIN;

CREATE DATABASE admin_api;

COMMIT;

BEGIN;

CREATE TABLE admin_api.orgs (
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

CREATE TABLE admin_api.failed_event (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE admin_api.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ,
    object_type TEXT,

    PRIMARY KEY (object_type)
);

CREATE TABLE admin_api.current_sequences (
    view_name TEXT,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);


COMMIT;
