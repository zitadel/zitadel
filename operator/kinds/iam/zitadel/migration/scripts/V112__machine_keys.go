package scripts

const V112MachineKeys = `
CREATE TABLE auth.machine_keys (
    id TEXT,
    user_id TEXT,

    machine_type SMALLINT,
    expiration_date TIMESTAMPTZ,
    sequence BIGINT,
    creation_date TIMESTAMPTZ,
    public_key JSONB,

    PRIMARY KEY (id, user_id)
);

ALTER TABLE management.machine_keys ADD COLUMN public_key JSONB;
`
