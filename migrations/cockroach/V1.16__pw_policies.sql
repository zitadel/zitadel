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