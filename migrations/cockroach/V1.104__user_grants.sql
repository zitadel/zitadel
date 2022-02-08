CREATE TABLE zitadel.projections.user_grants (
    id STRING NOT NULL
    , resource_owner STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL
    
    , user_id STRING NOT NULL
    , project_id STRING NOT NULL
    , grant_id STRING NOT NULL
    , roles STRING[]
    , state INT2 NOT NULL

    , PRIMARY KEY (id)
    , INDEX idx_user (user_id)
    , INDEX idx_ro (resource_owner)
);