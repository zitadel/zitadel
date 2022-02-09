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
    ,from STRING NOT NULL
    ,token JSONB

    , PRIMARY KEY (sms_id)
);