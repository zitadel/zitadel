CREATE TABLE zitadel.verifications(
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL

    , value TEXT
    , code BYTES

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expiry INTERVAL

    , failed_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_attempts >= 0)

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION zitadel.cleanup_verification(instance_id TEXT, verification_id TEXT)
    RETURNS VOID AS $$
BEGIN
    DELETE FROM zitadel.verifications v WHERE v.instance_id = instance_id AND v.id = verification_id;
END;
$$ LANGUAGE plpgsql;