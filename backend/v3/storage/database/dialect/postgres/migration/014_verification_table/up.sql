CREATE TABLE zitadel.verifications(
    instance_id TEXT NOT NULL
    , user_id TEXT
    , id TEXT NOT NULL DEFAULT gen_random_uuid()::TEXT

    , value TEXT
    , code BYTEA

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expiry INTERVAL

    , failed_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_attempts >= 0)

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

ALTER TABLE zitadel.users ADD CONSTRAINT fk_unverified_password FOREIGN KEY (instance_id, password_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (password_verification_id);
ALTER TABLE zitadel.users ADD CONSTRAINT fk_unverified_email FOREIGN KEY (instance_id, email_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (email_verification_id);
ALTER TABLE zitadel.users ADD CONSTRAINT fk_unverified_phone FOREIGN KEY (instance_id, phone_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (phone_verification_id);