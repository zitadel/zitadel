CREATE TABLE zitadel.projections.features
(
    aggregate_id                TEXT NOT NULL,
    creation_date               TIMESTAMPTZ NULL,
    change_date                 TIMESTAMPTZ NULL,
    sequence                    INT8 NULL,
    is_default                  BOOLEAN,

    tier_name                   TEXT,
    tier_description            TEXT,
    state                       INT2 NULL,
    state_description           TEXT,
    audit_log_retention         BIGINT,
    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,
    login_policy_password_reset BOOLEAN,
    password_complexity_policy  BOOLEAN,
    label_policy_private_label  BOOLEAN,
    label_policy_watermark      BOOLEAN,
    custom_domain               BOOLEAN,
    privacy_policy              BOOLEAN,
    meta_data_user              BOOLEAN,
    custom_text_message         BOOLEAN,
    custom_text_login           BOOLEAN,
    lockout_policy              BOOLEAN,
    actions                     BOOLEAN,

    PRIMARY KEY (aggregate_id)
);