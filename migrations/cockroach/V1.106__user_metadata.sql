CREATE TABLE zitadel.projections.user_metadatas (
    user_id STRING NOT NULL
    , resource_owner STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    
    , key STRING NOT NULL
    , value BYTEA

    , PRIMARY KEY (user_id, key)
);
