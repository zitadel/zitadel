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
    , sender_number STRING NOT NULL
    , host STRING NOT NULL
    , username STRING NOT NULL DEFAULT ''
    , password JSONB

    , PRIMARY KEY (aggregate_id)
);

CREATE TABLE zitadel.projections.sms_configs (
     id STRING NOT NULL
    ,aggregate_id STRING NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , resource_owner STRING NOT NULL
    , sequence INT8 NOT NULL
    , state INT2

    , PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.sms_configs_twilio (
    sms_id STRING NOT NULL
    ,sid STRING NOT NULL
    ,sender_name STRING NOT NULL
    ,token JSONB

    , PRIMARY KEY (sms_id)
);

ALTER TABLE zitadel.projections.login_policies
    ADD COLUMN password_check_lifetime BIGINT NOT NULL DEFAULT 0
    , ADD COLUMN external_login_check_lifetime BIGINT NOT NULL DEFAULT 0
    , ADD COLUMN mfa_init_skip_lifetime BIGINT NOT NULL DEFAULT 0
    , ADD COLUMN second_factor_check_lifetime BIGINT NOT NULL DEFAULT 0;
    , ADD COLUMN multi_factor_check_lifetime BIGINT NOT NULL DEFAULT 0;
