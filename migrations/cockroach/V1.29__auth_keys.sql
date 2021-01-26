-- ALTER TABLE auth.machine_keys
--     ADD COLUMN object_type STRING;
-- ALTER TABLE management.machine_key
--     ADD COLUMN object_type STRING;
--
-- ALTER TABLE auth.machine_keys
--     RENAME COLUMN user_id TO object_id;
-- ALTER TABLE management.machine_keys
--     RENAME COLUMN user_id TO object_id;
--
-- BEGIN;
--
-- ALTER TABLE auth.machine_keys
--     DROP CONSTRAINT "primary";
-- ALTER TABLE management.machine_keys
--     DROP CONSTRAINT "primary";
--
-- ALTER TABLE auth.machine_keys
--     ADD CONSTRAINT "primary" PRIMARY KEY (id, object_id, object_type);
-- ALTER TABLE management.machine_keys
--     ADD CONSTRAINT "primary" PRIMARY KEY (id, object_id, object_type);
--
-- COMMIT;

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