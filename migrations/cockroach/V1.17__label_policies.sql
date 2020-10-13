CREATE TABLE adminapi.label_policies (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    label_policy_state SMALLINT,
    sequence BIGINT,

    primary_color TEXT,
    secondary_color TEXT,

    PRIMARY KEY (aggregate_id)
);


CREATE TABLE management.label_policies (
      aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    label_policy_state SMALLINT,
    sequence BIGINT,

    primary_color TEXT,
    secondary_color TEXT,

    PRIMARY KEY (aggregate_id)
);

GRANT SELECT ON TABLE adminapi.label_policies TO notification;
GRANT SELECT ON TABLE management.label_policies TO notification;