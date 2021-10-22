CREATE TABLE zitadel.projections.login_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    allow_register BOOLEAN,
    allow_username_password BOOLEAN,
    allow_external_idps BOOLEAN,

    force_mfa BOOLEAN,
    second_factors SMALLINT ARRAY,
    multi_factors SMALLINT ARRAY,
    passwordless_type SMALLINT,
    is_default BOOLEAN,
    hide_password_reset BOOLEAN,

    PRIMARY KEY (aggregate_id)
);