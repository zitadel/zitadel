ALTER TABLE zitadel.users DROP CONSTRAINT IF EXISTS fk_unverified_password;
ALTER TABLE zitadel.users DROP CONSTRAINT IF EXISTS fk_unverified_email;
ALTER TABLE zitadel.users DROP CONSTRAINT IF EXISTS fk_unverified_phone;

DROP TRIGGER IF EXISTS user_verification_integrity_trigger ON zitadel.users;
DROP FUNCTION IF EXISTS zitadel.ensure_user_verification_integrity() CASCADE;
DROP TABLE IF EXISTS zitadel.verifications CASCADE;
