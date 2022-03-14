CREATE SCHEMA authz;

CREATE TABLE authz.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,

    PRIMARY KEY (projection_name)
);

CREATE TABLE authz.current_sequences (
    projection_name TEXT,
    aggregate_type TEXT,
    current_sequence BIGINT,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type)
);

CREATE TABLE authz.failed_events (
    projection_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,

    PRIMARY KEY (projection_name, failed_sequence)
);

CREATE TABLE authz.user_grants (
    id STRING NOT NULL,
    resource_owner STRING NULL,
    project_id STRING NULL,
    user_id STRING NULL,
    org_name STRING NULL,
    project_name STRING NULL,
    user_name STRING NULL,
    first_name STRING NULL,
    last_name STRING NULL,
    display_name STRING NULL,
    email STRING NULL,
    role_keys STRING[] NULL,
    grant_id STRING NULL,
    grant_state INT2 NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    org_primary_domain STRING NULL,
    project_owner STRING NULL,
    avatar_key STRING NULL,
    user_resource_owner STRING NULL,

    PRIMARY KEY (id)
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

    PRIMARY KEY (user_id, member_type, aggregate_id, object_id)
);
