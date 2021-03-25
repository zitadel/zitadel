CREATE TABLE adminapi.features
(
    aggregate_id                TEXT,

    creation_date               TIMESTAMPTZ,
    change_date                 TIMESTAMPTZ,
    sequence                    BIGINT,
    default_features            BOOLEAN,

    tier_name                   TEXT,
    tier_description            TEXT,
    state                       SMALLINT,
    state_description           TEXT,

    audit_log_retention         BIGINT,
    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,
    password_complexity_policy  BOOLEAN,
    label_policy                BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE auth.features
(
    aggregate_id                TEXT,

    creation_date               TIMESTAMPTZ,
    change_date                 TIMESTAMPTZ,
    sequence                    BIGINT,
    default_features            BOOLEAN,

    tier_name                   TEXT,
    tier_description            TEXT,
    state                       SMALLINT,
    state_description           TEXT,

    audit_log_retention         BIGINT,
    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,
    password_complexity_policy  BOOLEAN,
    label_policy                BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE authz.features
(
    aggregate_id                TEXT,

    creation_date               TIMESTAMPTZ,
    change_date                 TIMESTAMPTZ,
    sequence                    BIGINT,
    default_features            BOOLEAN,

    tier_name                   TEXT,
    tier_description            TEXT,
    state                       SMALLINT,
    state_description           TEXT,

    audit_log_retention         BIGINT,
    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,
    password_complexity_policy  BOOLEAN,
    label_policy                BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE management.features
(
    aggregate_id                TEXT,

    creation_date               TIMESTAMPTZ,
    change_date                 TIMESTAMPTZ,
    sequence                    BIGINT,
    default_features            BOOLEAN,

    tier_name                   TEXT,
    tier_description            TEXT,
    state                       SMALLINT,
    state_description           TEXT,

    audit_log_retention         BIGINT,
    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,
    password_complexity_policy  BOOLEAN,
    label_policy                BOOLEAN,

    PRIMARY KEY (aggregate_id)
);