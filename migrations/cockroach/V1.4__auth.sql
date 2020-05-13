BEGIN;

CREATE DATABASE auth;


COMMIT;

BEGIN;

CREATE USER auth;

GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE auth TO auth;
GRANT SELECT, INSERT, UPDATE ON DATABASE eventstore TO auth;
GRANT SELECT, INSERT, UPDATE ON TABLE eventstore.* TO auth;

COMMIT;

BEGIN;

CREATE TABLE auth.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ,
    object_type TEXT,

    PRIMARY KEY (object_type)
);

CREATE TABLE auth.current_sequences (
    view_name TEXT,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE auth.failed_event (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE auth.auth_requests (
    id TEXT,
    request JSONB,

    PRIMARY KEY (id)
);

CREATE TABLE auth.users (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    user_state SMALLINT,
    password_set BOOLEAN,
    password_change_required BOOLEAN,
    password_change TIMESTAMPTZ,
    last_login TIMESTAMPTZ,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nick_name TEXT,
    display_name TEXT,
    preferred_language TEXT,
    gender SMALLINT,
    email TEXT,
    is_email_verified BOOLEAN,
    phone TEXT,
    is_phone_verified BOOLEAN,
    country TEXT,
    locality TEXT,
    postal_code TEXT,
    region TEXT,
    street_address TEXT,
    otp_state SMALLINT,
    mfa_max_set_up SMALLINT,
    mfa_init_skipped TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

CREATE TABLE auth.user_sessions (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
--     state
    application_id TEXT,
    user_agent_id TEXT,
    user_id TEXT,
    user_name TEXT,
    password_verification TIMESTAMPTZ,
    mfa_software_verification TIMESTAMPTZ,
    mfa_hardware_verification TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

COMMIT;