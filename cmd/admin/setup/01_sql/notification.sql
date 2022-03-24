CREATE SCHEMA notification;

CREATE TABLE notification.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE notification.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    event_timestamp TIMESTAMPTZ,
    last_successful_spooler_run TIMESTAMPTZ,

    PRIMARY KEY (view_name)
);

CREATE TABLE notification.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE notification.notify_users (
    id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner STRING NULL,
    user_name STRING NULL,
    first_name STRING NULL,
    last_name STRING NULL,
    nick_name STRING NULL,
    display_name STRING NULL,
    preferred_language STRING NULL,
    gender INT2 NULL,
    last_email STRING NULL,
    verified_email STRING NULL,
    last_phone STRING NULL,
    verified_phone STRING NULL,
    sequence INT8 NULL,
    password_set BOOL NULL,
    login_names STRING NULL,
    preferred_login_name STRING NULL,
    instance_id STRING NULL,

    PRIMARY KEY (id)
);
