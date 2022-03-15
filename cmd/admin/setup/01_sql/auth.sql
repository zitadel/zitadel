CREATE SCHEMA auth;

CREATE TABLE auth.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,

    PRIMARY KEY (projection_name)
);

CREATE TABLE auth.current_sequences (
    projection_name TEXT,
    aggregate_type TEXT,
    current_sequence BIGINT,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type)
);

CREATE TABLE auth.failed_events (
    projection_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,

    PRIMARY KEY (projection_name, failed_sequence)
);

CREATE TABLE auth.users (
    id STRING NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner STRING NULL,
    user_state INT2 NULL,
    password_set BOOL NULL,
    password_change_required BOOL NULL,
    password_change TIMESTAMPTZ NULL,
    last_login TIMESTAMPTZ NULL,
    user_name STRING NULL,
    login_names STRING[] NULL,
    preferred_login_name STRING NULL,
    first_name STRING NULL,
    last_name STRING NULL,
    nick_name STRING NULL,
    display_name STRING NULL,
    preferred_language STRING NULL,
    gender INT2 NULL,
    email STRING NULL,
    is_email_verified BOOL NULL,
    phone STRING NULL,
    is_phone_verified BOOL NULL,
    country STRING NULL,
    locality STRING NULL,
    postal_code STRING NULL,
    region STRING NULL,
    street_address STRING NULL,
    otp_state INT2 NULL,
    mfa_max_set_up INT2 NULL,
    mfa_init_skipped TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    init_required BOOL NULL,
    username_change_required BOOL NULL,
    machine_name STRING NULL,
    machine_description STRING NULL,
    user_type STRING NULL,
    u2f_tokens BYTES NULL,
    passwordless_tokens BYTES NULL,
    avatar_key STRING NULL,
    passwordless_init_required BOOL NULL,
    password_init_required BOOL NULL,
    tenant STRING NULL,

    PRIMARY KEY (id)
);

CREATE TABLE auth.user_sessions (
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner STRING NULL,
    state INT2 NULL,
    user_agent_id STRING NULL,
    user_id STRING NULL,
    user_name STRING NULL,
    password_verification TIMESTAMPTZ NULL,
    second_factor_verification TIMESTAMPTZ NULL,
    multi_factor_verification TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    second_factor_verification_type INT2 NULL,
    multi_factor_verification_type INT2 NULL,
    user_display_name STRING NULL,
    login_name STRING NULL,
    external_login_verification TIMESTAMPTZ NULL,
    selected_idp_config_id STRING NULL,
    passwordless_verification TIMESTAMPTZ NULL,
    avatar_key STRING NULL,
    tenant STRING NULL,

    PRIMARY KEY (user_agent_id, user_id)
);

CREATE TABLE auth.user_external_idps (
    external_user_id STRING NOT NULL,
    idp_config_id STRING NOT NULL,
    user_id STRING NULL,
    idp_name STRING NULL,
    user_display_name STRING NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    resource_owner STRING NULL,
    tenant STRING NULL,

    PRIMARY KEY (external_user_id, idp_config_id)
);

CREATE TABLE auth.tokens (
    id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner STRING NULL,
    application_id STRING NULL,
    user_agent_id STRING NULL,
    user_id STRING NULL,
    expiration TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    scopes STRING[] NULL,
    audience STRING[] NULL,
    preferred_language STRING NULL,
    refresh_token_id STRING NULL,
    tenant STRING NULL,

    PRIMARY KEY (id),
    INDEX user_user_agent_idx (user_id, user_agent_id)
);

CREATE TABLE auth.refresh_tokens (
    id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    resource_owner STRING NULL,
    token STRING NULL,
    client_id STRING NOT NULL,
    user_agent_id STRING NOT NULL,
    user_id STRING NOT NULL,
    auth_time TIMESTAMPTZ NULL,
    idle_expiration TIMESTAMPTZ NULL,
    expiration TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    scopes STRING[] NULL,
    audience STRING[] NULL,
    amr STRING[] NULL,
    tenant STRING NULL,

    PRIMARY KEY (id),
    UNIQUE INDEX unique_client_user_index (client_id ASC, user_agent_id ASC, user_id ASC)
);

CREATE TABLE auth.org_project_mapping (
    org_id STRING NOT NULL,
    project_id STRING NOT NULL,
    project_grant_id STRING NULL,
    tenant STRING NULL,

    PRIMARY KEY (org_id, project_id)
);

CREATE TABLE auth.idp_providers (
    aggregate_id STRING NOT NULL,
    idp_config_id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    name STRING NULL,
    idp_config_type INT2 NULL,
    idp_provider_type INT2 NULL,
    idp_state INT2 NULL,
    styling_type INT2 NULL,
    tenant STRING NULL,

    PRIMARY KEY (aggregate_id, idp_config_id)
);

CREATE TABLE auth.idp_configs (
    idp_config_id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    aggregate_id STRING NULL,
    name STRING NULL,
    idp_state INT2 NULL,
    idp_provider_type INT2 NULL,
    is_oidc BOOL NULL,
    oidc_client_id STRING NULL,
    oidc_client_secret JSONB NULL,
    oidc_issuer STRING NULL,
    oidc_scopes STRING[] NULL,
    oidc_idp_display_name_mapping INT2 NULL,
    oidc_idp_username_mapping INT2 NULL,
    styling_type INT2 NULL,
    oauth_authorization_endpoint STRING NULL,
    oauth_token_endpoint STRING NULL,
    auto_register BOOL NULL,
    jwt_endpoint STRING NULL,
    jwt_keys_endpoint STRING NULL,
    jwt_header_name STRING NULL,
    tenant STRING NULL,

    PRIMARY KEY (idp_config_id)
);

CREATE TABLE auth.auth_requests (
    id STRING NOT NULL,
    request JSONB NULL,
    code STRING NULL,
    request_type INT2 NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    tenant STRING NULL,

    PRIMARY KEY (id),
    INDEX auth_code_idx (code)
);
