CREATE TABLE management.machine_keys (
    id TEXT,
    user_id TEXT,

    machine_type SMALLINT,
    expiration_date TIMESTAMPTZ,
    sequence BIGINT,
    creation_date TIMESTAMPTZ,

    PRIMARY KEY (id, user_id)
)