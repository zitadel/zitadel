CREATE TABLE zitadel.projections.authn_keys(
    id STRING
    , creation_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , aggregate_id STRING NOT NULL
    , sequence INT8 NOT NULL

    , object_id STRING NOT NULL
    , expiration TIMESTAMPTZ NOT NULL
    , identifier STRING NOT NULL
    , public_key BYTES NOT NULL
    , enabled BOOLEAN NOT NULL DEFAULT true

    , PRIMARY KEY (id)
);
