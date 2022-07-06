CREATE SCHEMA notification;

CREATE TABLE notification.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, instance_id)
);

CREATE TABLE notification.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    event_timestamp TIMESTAMPTZ,
    last_successful_spooler_run TIMESTAMPTZ,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, instance_id)
);

CREATE TABLE notification.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, failed_sequence, instance_id)
);

CREATE TABLE notification.notify_users (
    id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner TEXT NULL,
    user_name TEXT NULL,
    first_name TEXT NULL,
    last_name TEXT NULL,
    nick_name TEXT NULL,
    display_name TEXT NULL,
    preferred_language TEXT NULL,
    gender INT2 NULL,
    last_email TEXT NULL,
    verified_email TEXT NULL,
    last_phone TEXT NULL,
    verified_phone TEXT NULL,
    sequence INT8 NULL,
    password_set BOOL NULL,
    login_names TEXT NULL,
    preferred_login_name TEXT NULL,
    instance_id TEXT NULL,

    PRIMARY KEY (id)
);
