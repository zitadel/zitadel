
ALTER TABLE adminapi.idp_configs ADD COLUMN oidc_idp_display_name_mapping SMALLINT;
ALTER TABLE adminapi.idp_configs ADD COLUMN oidc_idp_username_mapping SMALLINT;

ALTER TABLE management.idp_configs ADD COLUMN oidc_idp_display_name_mapping SMALLINT;
ALTER TABLE management.idp_configs ADD COLUMN oidc_idp_username_mapping SMALLINT;

CREATE TABLE auth.idp_configs (
    idp_config_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    aggregate_id TEXT,
    name TEXT,
    logo_src BYTES,
    idp_state SMALLINT,
    idp_provider_type SMALLINT,

    is_oidc BOOLEAN,
    oidc_client_id TEXT,
    oidc_client_secret JSONB,
    oidc_issuer TEXT,
    oidc_scopes TEXT ARRAY,
    oidc_idp_display_name_mapping SMALLINT,
    oidc_idp_username_mapping SMALLINT,

    PRIMARY KEY (idp_config_id)
);

CREATE TABLE auth.login_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    login_policy_state SMALLINT,
    sequence BIGINT,

    allow_register BOOLEAN,
    allow_username_password BOOLEAN,
    allow_external_idp BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE auth.idp_providers (
    aggregate_id TEXT,
    idp_config_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    name string,
    idp_config_type SMALLINT,
    idp_provider_type SMALLINT,

    PRIMARY KEY (aggregate_id, idp_config_id)
);

CREATE TABLE auth.user_external_idps (
    external_user_id TEXT,
    idp_config_id TEXT,
    user_id TEXT,
    idp_name TEXT,
    user_display_name TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    resource_owner TEXT,

    PRIMARY KEY (external_user_id, idp_config_id)
);

CREATE TABLE management.user_external_idps (
    idp_config_id TEXT,
    external_user_id TEXT,
    user_id TEXT,
    idp_name TEXT,
    user_display_name TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    resource_owner TEXT,

    PRIMARY KEY (external_user_id, idp_config_id)
);
 
CREATE TABLE adminapi.user_external_idps (
    idp_config_id TEXT,
    external_user_id TEXT,
    user_id TEXT,
    idp_name TEXT,
    user_display_name TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    resource_owner TEXT,

    PRIMARY KEY (external_user_id, idp_config_id)
);
