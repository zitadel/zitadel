CREATE TABLE zitadel.projections.keys(
    id STRING,
    is_private BOOLEAN,
    creation_date TIMESTAMPTZ NOT NULL,
    change_date TIMESTAMPTZ NOT NULL,
    resource_owner STRING NOT NULL,
    sequence INT8 NOT NULL,
    algorithm STRING DEFAULT '' NOT NULL,
    use STRING DEFAULT '' NOT NULL,
    expiry TIMESTAMPTZ NOT NULL,
    key JSONB NOT NULL,

    PRIMARY KEY (id, is_private)
);
