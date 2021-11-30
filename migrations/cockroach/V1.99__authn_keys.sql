CREATE TABLE zitadel.projections.authn_keys(
    id STRING
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , aggregate_id STRING NOT NULL
    , sequence INT8 NOT NULL

    , object_id STRING NOT NULL
    , expiration TIMESTAMPTZ NOT NULL

    , PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.authn_keys_authentication(
	key_id STRING REFERENCES zitadel.projections.authn_keys (id) ON DELETE CASCADE
	
    , identifier STRING NOT NULL
    , key BYTES NOT NULL

    , PRIMARY KEY (key_id)
);
