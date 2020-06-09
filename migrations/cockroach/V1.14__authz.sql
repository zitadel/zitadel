BEGIN;

CREATE DATABASE authz;


COMMIT;

BEGIN;

CREATE USER authz;

GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE authz TO authz;
GRANT SELECT, INSERT, UPDATE ON DATABASE eventstore TO authz;
GRANT SELECT, INSERT, UPDATE ON TABLE eventstore.* TO authz;
GRANT SELECT, INSERT, UPDATE ON DATABASE auth TO authz;
GRANT SELECT, INSERT, UPDATE ON TABLE auth.* TO authz;

COMMIT;

BEGIN;

CREATE TABLE authz.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ,
    object_type TEXT,

    PRIMARY KEY (object_type)
);

CREATE TABLE authz.current_sequences (
    view_name TEXT,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE authz.failed_event (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE authz.user_grants (
    id TEXT,
    resource_owner TEXT,
    project_id TEXT,
    user_id TEXT,
    org_name TEXT,
    org_domain TEXT,
    project_name TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    role_keys TEXT Array,

    grant_state SMALLINT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

CREATE TABLE authz.applications (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    app_state SMALLINT,
    resource_owner TEXT,
    app_name TEXT,
    project_id TEXT,
    app_type SMALLINT,
    is_oidc BOOLEAN,
    oidc_client_id TEXT,
    oidc_redirect_uris TEXT ARRAY,
    oidc_response_types SMALLINT ARRAY,
    oidc_grant_types SMALLINT ARRAY,
    oidc_application_type SMALLINT,
    oidc_auth_method_type SMALLINT,
    oidc_post_logout_redirect_uris TEXT ARRAY,

    PRIMARY KEY (id)
);

COMMIT;

