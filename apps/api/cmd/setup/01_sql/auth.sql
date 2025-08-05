CREATE SCHEMA auth;

CREATE TABLE auth.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, instance_id)
);

CREATE TABLE auth.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    event_date TIMESTAMPTZ,
    last_successful_spooler_run TIMESTAMPTZ,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, instance_id)
);

CREATE TABLE auth.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, failed_sequence, instance_id)
);

CREATE TABLE auth.users (
    id TEXT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner TEXT NULL,
    user_state INT2 NULL,
    password_set BOOL NULL,
    password_change_required BOOL NULL,
    password_change TIMESTAMPTZ NULL,
    last_login TIMESTAMPTZ NULL,
    user_name TEXT NULL,
    login_names TEXT[] NULL,
    preferred_login_name TEXT NULL,
    first_name TEXT NULL,
    last_name TEXT NULL,
    nick_name TEXT NULL,
    display_name TEXT NULL,
    preferred_language TEXT NULL,
    gender INT2 NULL,
    email TEXT NULL,
    is_email_verified BOOL NULL,
    phone TEXT NULL,
    is_phone_verified BOOL NULL,
    country TEXT NULL,
    locality TEXT NULL,
    postal_code TEXT NULL,
    region TEXT NULL,
    street_address TEXT NULL,
    otp_state INT2 NULL,
    mfa_max_set_up INT2 NULL,
    mfa_init_skipped TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    init_required BOOL NULL,
    username_change_required BOOL NULL,
    machine_name TEXT NULL,
    machine_description TEXT NULL,
    user_type TEXT NULL,
    u2f_tokens BYTEA NULL,
    passwordless_tokens BYTEA NULL,
    avatar_key TEXT NULL,
    passwordless_init_required BOOL NULL,
    password_init_required BOOL NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (id, instance_id)
);

CREATE TABLE auth.user_sessions (
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner TEXT NULL,
    state INT2 NULL,
    user_agent_id TEXT NULL,
    user_id TEXT NULL,
    user_name TEXT NULL,
    password_verification TIMESTAMPTZ NULL,
    second_factor_verification TIMESTAMPTZ NULL,
    multi_factor_verification TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    second_factor_verification_type INT2 NULL,
    multi_factor_verification_type INT2 NULL,
    user_display_name TEXT NULL,
    login_name TEXT NULL,
    external_login_verification TIMESTAMPTZ NULL,
    selected_idp_config_id TEXT NULL,
    passwordless_verification TIMESTAMPTZ NULL,
    avatar_key TEXT NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (user_agent_id, user_id, instance_id)
);

CREATE TABLE auth.user_external_idps (
    external_user_id TEXT NOT NULL,
    idp_config_id TEXT NOT NULL,
    user_id TEXT NULL,
    idp_name TEXT NULL,
    user_display_name TEXT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    resource_owner TEXT NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (external_user_id, idp_config_id, instance_id)
);

CREATE TABLE auth.tokens (
    id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner TEXT NULL,
    application_id TEXT NULL,
    user_agent_id TEXT NULL,
    user_id TEXT NULL,
    expiration TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    scopes TEXT[] NULL,
    audience TEXT[] NULL,
    preferred_language TEXT NULL,
    refresh_token_id TEXT NULL,
    is_pat BOOL NOT NULL DEFAULT false,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (id, instance_id)
);

CREATE INDEX user_user_agent_idx ON auth.tokens (user_id, user_agent_id);

CREATE TABLE auth.refresh_tokens (
    id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner TEXT NULL,
    token TEXT NULL,
    client_id TEXT NOT NULL,
    user_agent_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    auth_time TIMESTAMPTZ NULL,
    idle_expiration TIMESTAMPTZ NULL,
    expiration TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    scopes TEXT[] NULL,
    audience TEXT[] NULL,
    amr TEXT[] NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (id, instance_id)
);

CREATE UNIQUE INDEX unique_client_user_index ON auth.refresh_tokens (client_id, user_agent_id, user_id);

CREATE TABLE auth.org_project_mapping (
    org_id TEXT NOT NULL,
    project_id TEXT NOT NULL,
    project_grant_id TEXT NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (org_id, project_id, instance_id)
);

CREATE TABLE auth.idp_providers (
    aggregate_id TEXT NOT NULL,
    idp_config_id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    name TEXT NULL,
    idp_config_type INT2 NULL,
    idp_provider_type INT2 NULL,
    idp_state INT2 NULL,
    styling_type INT2 NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (aggregate_id, idp_config_id, instance_id)
);

CREATE TABLE auth.idp_configs (
    idp_config_id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    aggregate_id TEXT NULL,
    name TEXT NULL,
    idp_state INT2 NULL,
    idp_provider_type INT2 NULL,
    is_oidc BOOL NULL,
    oidc_client_id TEXT NULL,
    oidc_client_secret JSONB NULL,
    oidc_issuer TEXT NULL,
    oidc_scopes TEXT[] NULL,
    oidc_idp_display_name_mapping INT2 NULL,
    oidc_idp_username_mapping INT2 NULL,
    styling_type INT2 NULL,
    oauth_authorization_endpoint TEXT NULL,
    oauth_token_endpoint TEXT NULL,
    auto_register BOOL NULL,
    jwt_endpoint TEXT NULL,
    jwt_keys_endpoint TEXT NULL,
    jwt_header_name TEXT NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (idp_config_id, instance_id)
);

CREATE TABLE auth.auth_requests (
    id TEXT NOT NULL,
    request JSONB NULL,
    code TEXT NULL,
    request_type INT2 NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (id, instance_id)
);

CREATE INDEX auth_code_idx ON auth.auth_requests (code);
