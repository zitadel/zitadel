CREATE TABLE zitadel.user_verifications(
    instance_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    , id TEXT NOT NULL DEFAULT gen_random_uuid()::TEXT

    , value TEXT
    , code BYTEA

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expiry INTERVAL

    , failed_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_attempts >= 0)

    , PRIMARY KEY (instance_id, user_id, id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, user_id) ON DELETE CASCADE
);
