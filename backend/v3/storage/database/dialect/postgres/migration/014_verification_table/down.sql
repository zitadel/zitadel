ALTER TABLE zitadel.users DROP CONSTRAINT fk_unverified_password;
ALTER TABLE zitadel.users DROP CONSTRAINT fk_unverified_email;
ALTER TABLE zitadel.users DROP CONSTRAINT fk_unverified_phone;
ALTER TABLE zitadel.users DROP CONSTRAINT fk_email_otp_verification;
ALTER TABLE zitadel.users DROP CONSTRAINT fk_sms_otp_verification;
ALTER TABLE zitadel.human_passkeys DROP CONSTRAINT fk_init_verification

DROP FUNCTION IF EXISTS zitadel.cleanup_verification(TEXT, TEXT);
DROP TABLE IF EXISTS zitadel.verifications;
