
ALTER TABLE management.users ADD COLUMN machine_name STRING, ADD COLUMN machine_description STRING, ADD COLUMN user_type STRING;
ALTER TABLE adminapi.users ADD COLUMN machine_name STRING, ADD COLUMN machine_description STRING, ADD COLUMN user_type STRING;
ALTER TABLE auth.users ADD COLUMN machine_name STRING, ADD COLUMN machine_description STRING, ADD COLUMN user_type STRING;

CREATE TABLE management.machine_keys (
    id TEXT,
    user_id TEXT,

    machine_type SMALLINT,
    expiration_date TIMESTAMPTZ,
    sequence BIGINT,
    creation_date TIMESTAMPTZ,

    PRIMARY KEY (id, user_id)
)
