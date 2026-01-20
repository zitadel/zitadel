CREATE TABLE zitadel.verifications(
    instance_id TEXT NOT NULL
    , user_id TEXT
    , id TEXT NOT NULL DEFAULT gen_random_uuid()::TEXT

    , value TEXT
    , code BYTEA

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expiry INTERVAL

    , failed_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_attempts >= 0)

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

ALTER TABLE zitadel.users ADD CONSTRAINT fk_unverified_password FOREIGN KEY (instance_id, unverified_password_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_password_id);
ALTER TABLE zitadel.users ADD CONSTRAINT fk_unverified_email FOREIGN KEY (instance_id, unverified_email_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_email_id);
ALTER TABLE zitadel.users ADD CONSTRAINT fk_unverified_phone FOREIGN KEY (instance_id, unverified_phone_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_phone_id);
ALTER TABLE zitadel.users ADD CONSTRAINT fk_email_otp_verification FOREIGN KEY (instance_id, email_otp_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (email_otp_verification_id);
ALTER TABLE zitadel.users ADD CONSTRAINT fk_sms_otp_verification FOREIGN KEY (instance_id, sms_otp_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (sms_otp_verification_id);
ALTER TABLE zitadel.human_passkeys ADD CONSTRAINT fk_init_verification FOREIGN KEY (instance_id, init_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (init_verification_id);