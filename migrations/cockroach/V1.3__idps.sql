
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
