CREATE TABLE zitadel.verifications(
    instance_id TEXT NOT NULL
    , user_id TEXT
    , id TEXT NOT NULL DEFAULT gen_random_uuid()::TEXT

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

CREATE OR REPLACE FUNCTION zitadel.ensure_user_verification_integrity() RETURNS trigger AS $$
BEGIN
    IF (NEW.password_verification_id IS DISTINCT FROM OLD.password_verification_id) THEN
        DELETE FROM zitadel.verifications
        WHERE instance_id = NEW.instance_id
        AND id = OLD.password_verification_id;
    END IF;

    IF (NEW.email_verification_id IS DISTINCT FROM OLD.email_verification_id) THEN
        DELETE FROM zitadel.verifications
        WHERE instance_id = NEW.instance_id
        AND id = OLD.email_verification_id;
    END IF;

    IF (NEW.phone_verification_id IS DISTINCT FROM OLD.phone_verification_id) THEN
        DELETE FROM zitadel.verifications
        WHERE instance_id = NEW.instance_id
        AND id = OLD.phone_verification_id;
    END IF;
    
    IF (NEW.invite_verification_id IS DISTINCT FROM OLD.invite_verification_id) THEN
        DELETE FROM zitadel.verifications
        WHERE instance_id = NEW.instance_id
        AND id = OLD.invite_verification_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_verification_integrity_trigger
AFTER UPDATE ON zitadel.users
FOR EACH ROW
EXECUTE FUNCTION zitadel.ensure_user_verification_integrity();

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.verifications
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();