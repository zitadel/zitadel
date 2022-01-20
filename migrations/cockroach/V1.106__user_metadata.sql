CREATE TABLE zitadel.projections.user_metadata (
    user_id STRING NOT NULL
    , resource_owner STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    
    , key STRING NOT NULL
    , value BYTES

    , PRIMARY KEY (user_id, key)
    , INDEX idx_ro (resource_owner)
);
