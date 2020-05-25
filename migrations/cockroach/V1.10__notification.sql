BEGIN;

CREATE DATABASE notification;


COMMIT;

BEGIN;

CREATE USER notification;

GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE notification TO notification;
GRANT SELECT, INSERT, UPDATE ON DATABASE eventstore TO notification;
GRANT SELECT, INSERT, UPDATE ON TABLE eventstore.* TO notification;

COMMIT;

BEGIN;

CREATE TABLE notification.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ,
    object_type TEXT,

    PRIMARY KEY (object_type)
);

CREATE TABLE notification.current_sequences (
    view_name TEXT,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE notification.failed_event (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE notification.notify_users (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nick_Name TEXT,
    display_name TEXT,
    preferred_language TEXT,
    gender SMALLINT,
    last_email TEXT,
    verified_email TEXT,
    last_phone TEXT,
    verified_phone TEXT,
    sequence BIGINT,

    PRIMARY KEY (id)
);


COMMIT;
