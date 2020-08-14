
CREATE TABLE adminapi.idp_configs (
    idp_config_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    iam_id TEXT,
    name TEXT,
    logo_src TEXT,
    idp_state SMALLINT,

    is_oidc BOOLEAN,
    oidc_client_id TEXT,
    oidc_client_secret BYTEA,
    oidc_issuer TEXT,
    oidc_scopes TEXT ARRAY,

    PRIMARY KEY (idp_config_id)
);


CREATE TABLE management.idp_configs (
    idp_config_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    resource_owner TEXT,
    name TEXT,
    logo_src TEXT,
    idp_state SMALLINT,

    is_oidc BOOLEAN,
    oidc_client_id TEXT,
    oidc_client_secret BYTEA,
    oidc_issuer TEXT,
    oidc_scopes TEXT ARRAY,

    PRIMARY KEY (idp_config_id)
);


CREATE TABLE adminapi.login_policies (
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


CREATE TABLE management.login_policies (
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

CREATE TABLE adminapi.idp_providers (
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

CREATE TABLE management.idp_providers (
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