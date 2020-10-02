CREATE TABLE auth.password_complexity_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    complexity_policy_state SMALLINT,
    sequence BIGINT,

    min_length BIGINT,
    has_lowercase BOOLEAN,
    has_uppercase BOOLEAN,
    has_symbol BOOLEAN,
    has_number BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE adminapi.password_complexity_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    complexity_policy_state SMALLINT,
    sequence BIGINT,

    min_length BIGINT,
    has_lowercase BOOLEAN,
    has_uppercase BOOLEAN,
    has_symbol BOOLEAN,
    has_number BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE management.password_complexity_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    complexity_policy_state SMALLINT,
    sequence BIGINT,

    min_length BIGINT,
    has_lowercase BOOLEAN,
    has_uppercase BOOLEAN,
    has_symbol BOOLEAN,
    has_number BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE auth.password_age_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    age_policy_state SMALLINT,
    sequence BIGINT,

    max_age_days BIGINT,
    expire_warn_days BIGINT,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE adminapi.password_age_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    age_policy_state SMALLINT,
    sequence BIGINT,

    max_age_days BIGINT,
    expire_warn_days BIGINT,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE management.password_age_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    age_policy_state SMALLINT,
    sequence BIGINT,

    max_age_days BIGINT,
    expire_warn_days BIGINT,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE auth.password_lockout_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    lockout_policy_state SMALLINT,
    sequence BIGINT,

    max_attempts BIGINT,
    show_lockout_failures BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE adminapi.password_lockout_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    lockout_policy_state SMALLINT,
    sequence BIGINT,

    max_attempts BIGINT,
    show_lockout_failures BOOLEAN,

    PRIMARY KEY (aggregate_id)
);

CREATE TABLE management.password_lockout_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    lockout_policy_state SMALLINT,
    sequence BIGINT,

    max_attempts BIGINT,
    show_lockout_failures BOOLEAN,

    PRIMARY KEY (aggregate_id)
);