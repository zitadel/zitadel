CREATE TABLE zitadel.projections.user_auth_methods (
    token_id STRING NOT NULL
    , resource_owner STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL

    , user_id STRING NOT NULL
    , state INT2 NOT NULL
    , method_type INT2 NOT NULL
    , name STRING NOT NULL

    , PRIMARY KEY (user_id, method_type, token_id)
    , INDEX idx_ro (resource_owner)
);
