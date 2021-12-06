CREATE TABLE zitadel.projections.keys (
    id STRING,
    creation_date TIMESTAMPTZ NOT NULL,
    change_date TIMESTAMPTZ NOT NULL,
    resource_owner STRING NOT NULL,
    sequence INT8 NOT NULL,
    algorithm STRING DEFAULT '' NOT NULL,
    use STRING DEFAULT '' NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.keys_private(
    id STRING REFERENCES zitadel.projections.keys ON DELETE NO ACTION,
    expiry TIMESTAMPTZ NOT NULL,
    key JSONB NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.keys_public (
    id STRING,
    expiry TIMESTAMPTZ NOT NULL,
    key BYTES NOT NULL,

    PRIMARY KEY (id)
);
