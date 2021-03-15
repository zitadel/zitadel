CREATE TABLE adminapi.features
(
    aggregate_id                TEXT,

    creation_date               TIMESTAMPTZ,
    change_date                 TIMESTAMPTZ,
    sequence                    BIGINT,
    default_features            BOOLEAN,

    tier_name                   TEXT,
    tier_description            TEXT,
    tier_state                  SMALLINT,
    tier_state_description      TEXT,

    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,

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
    tier_state                  SMALLINT,
    tier_state_description      TEXT,

    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,

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
    tier_state                  SMALLINT,
    tier_state_description      TEXT,

    login_policy_factors        BOOLEAN,
    login_policy_idp            BOOLEAN,
    login_policy_passwordless   BOOLEAN,
    login_policy_registration   BOOLEAN,
    login_policy_username_login BOOLEAN,

    PRIMARY KEY (aggregate_id)
);