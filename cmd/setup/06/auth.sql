CREATE TABLE auth.users2 (
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
    owner_removed BOOL DEFAULT false,

    PRIMARY KEY (instance_id, id)
);
CREATE INDEX u2_owner_removed_idx ON auth.users2 (owner_removed);

CREATE TABLE auth.user_external_idps2 (
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
    owner_removed BOOL DEFAULT false,

    PRIMARY KEY (instance_id, external_user_id, idp_config_id)
);
CREATE INDEX ext_idps2_owner_removed_idx ON auth.user_external_idps2 (owner_removed);

CREATE TABLE auth.org_project_mapping2 (
    org_id TEXT NOT NULL,
    project_id TEXT NOT NULL,
    project_grant_id TEXT NULL,
    instance_id TEXT NOT NULL,
    owner_removed BOOL DEFAULT false,

    PRIMARY KEY (instance_id, org_id, project_id)
);
CREATE INDEX org_proj_m2_owner_removed_idx ON auth.org_project_mapping2 (owner_removed);

CREATE TABLE auth.idp_providers2 (
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
    owner_removed BOOL DEFAULT false,

    PRIMARY KEY (instance_id, aggregate_id, idp_config_id)
);
CREATE INDEX idp_prov2_owner_removed_idx ON auth.idp_providers2 (owner_removed);

CREATE TABLE auth.idp_configs2 (
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
    owner_removed BOOL DEFAULT false,

    PRIMARY KEY (instance_id, idp_config_id)
);
CREATE INDEX idp_conf2_owner_removed_idx ON auth.idp_configs2 (owner_removed);
