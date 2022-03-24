CREATE SCHEMA authz;

CREATE TABLE authz.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE authz.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    event_timestamp TIMESTAMPTZ,
    last_successful_spooler_run TIMESTAMPTZ,

    PRIMARY KEY (view_name)
);

CREATE TABLE authz.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE authz.user_memberships (
    user_id STRING NOT NULL,
    member_type INT2 NOT NULL,
    aggregate_id STRING NOT NULL,
    object_id STRING NOT NULL,
    roles STRING[] NULL,
    display_name STRING NULL,
    resource_owner STRING NULL,
    resource_owner_name STRING NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    instance_id STRING NULL,

    PRIMARY KEY (user_id, member_type, aggregate_id, object_id)
);
