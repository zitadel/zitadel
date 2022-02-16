ALTER TABLE zitadel.projections.iam ADD COLUMN default_language TEXT DEFAULT '';

CREATE TABLE zitadel.projections.secret_generators (
    generator_type INT2 NOT NULL
    , aggregate_id STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , sequence INT8 NOT NULL

    , length BIGINT NOT NULL
    , expiry BIGINT NOT NULL
    , include_lower_letters BOOLEAN NOT NULL
    , include_upper_letters BOOLEAN NOT NULL
    , include_digits BOOLEAN NOT NULL
    , include_symbols BOOLEAN NOT NULL

    , PRIMARY KEY (generator_type, aggregate_id)
);

CREATE TABLE zitadel.projections.smtp_configs (
    aggregate_id STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , sequence INT8 NOT NULL

    , tls BOOLEAN NOT NULL
    , sender_address STRING NOT NULL
    , sender_name STRING NOT NULL
    , host STRING NOT NULL
    , username STRING NOT NULL DEFAULT ''
    , password JSONB

    , PRIMARY KEY (aggregate_id)
);
