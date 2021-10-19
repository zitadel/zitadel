CREATE TABLE zitadel.projections.idps (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,
    resource_owner TEXT,
    
    state SMALLINT,
    name TEXT,
    styling_type SMALLINT,
    owner_type SMALLINT,
    auto_register BOOLEAN,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.idps_oidc_config (
    idp_id TEXT REFERENCES zitadel.projections.idps (id) ON DELETE CASCADE,
    
    client_id TEXT,
    client_secret JSONB,
    issuer TEXT,
    scopes STRING[],
    display_name_mapping SMALLINT,
    username_mapping SMALLINT,
    authorization_endpoint TEXT,
    token_endpoint TEXT,

    PRIMARY KEY (idp_id)
);

CREATE TABLE zitadel.projections.idps_jwt_config (
    idp_id TEXT REFERENCES zitadel.projections.idps (id) ON DELETE CASCADE,

    issuer TEXT,
    keys_endpoint TEXT,
    header_name TEXT,
    endpoint TEXT,

    PRIMARY KEY (idp_id)
);