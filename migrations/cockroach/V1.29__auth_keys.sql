CREATE TABLE auth.authn_keys
(
    key_id          TEXT,
    object_id       TEXT,
    object_type     SMALLINT,

    key_type        SMALLINT,
    sequence        BIGINT,
    expiration_date TIMESTAMPTZ,
    creation_date   TIMESTAMPTZ,
    public_key      BYTES,

    PRIMARY KEY (key_id, object_id, object_type)
);

INSERT INTO auth.authn_keys (
    key_id,
    object_id,
    object_type,
    key_type,
    sequence,
    expiration_date,
    creation_date,
    public_key
    )
    SELECT
        id,
        user_id,
        0,
        machine_type,
        sequence,
        expiration_date,
        creation_date,
        public_key
    FROM auth.machine_keys;

CREATE TABLE management.authn_keys
(
    key_id          TEXT,
    object_id       TEXT,
    object_type     SMALLINT,

    key_type        SMALLINT,
    sequence        BIGINT,
    expiration_date TIMESTAMPTZ,
    creation_date   TIMESTAMPTZ,
    public_key      BYTES,

    PRIMARY KEY (key_id, object_id, object_type)
);

INSERT INTO management.authn_keys (
    key_id,
    object_id,
    object_type,
    key_type,
    sequence,
    expiration_date,
    creation_date,
    public_key
)
SELECT
    id,
    user_id,
    0,
    machine_type,
    sequence,
    expiration_date,
    creation_date,
    public_key
FROM management.machine_keys;