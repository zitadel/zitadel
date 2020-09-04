
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
