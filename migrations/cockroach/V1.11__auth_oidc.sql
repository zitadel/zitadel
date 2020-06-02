BEGIN;

CREATE TABLE auth.keys (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    private BOOLEAN,
    expiry TIMESTAMPTZ,
    algorithm TEXT,
    usage SMALLINT,
    key JSONB,
    sequence BIGINT,

    PRIMARY KEY (id, private)
);

CREATE TABLE auth.applications (
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

ALTER TABLE auth.tokens ADD COLUMN scopes TEXT ARRAY;

ALTER TABLE auth.users ADD COLUMN init_required BOOLEAN;
ALTER TABLE management.users ADD COLUMN init_required BOOLEAN;

ALTER TABLE auth.auth_requests ADD COLUMN code TEXT;
ALTER TABLE auth.auth_requests ADD COLUMN request_type smallint;

ALTER TABLE auth.user_sessions ADD COLUMN mfa_software_verification_type smallint;
ALTER TABLE auth.user_sessions ADD COLUMN mfa_hardware_verification_type smallint;

COMMIT;