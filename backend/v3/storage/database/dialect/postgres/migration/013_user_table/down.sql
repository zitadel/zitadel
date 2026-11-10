DROP TABLE IF EXISTS zitadel.user_identity_provider_links;
DROP TABLE IF EXISTS zitadel.user_passkeys;
DROP TYPE IF EXISTS zitadel.passkey_type;
DROP TABLE IF EXISTS zitadel.machine_keys;
DROP TABLE IF EXISTS zitadel.user_personal_access_tokens;
DROP TABLE IF EXISTS zitadel.user_metadata CASCADE;

ALTER TABLE zitadel.identity_provider_intents DROP CONSTRAINT IF EXISTS fk_idp_intent_user;
ALTER TABLE zitadel.sessions DROP CONSTRAINT IF EXISTS fk_session_user;
ALTER TABLE zitadel.authorizations DROP CONSTRAINT IF EXISTS fk_authorization_user;

DROP TABLE IF EXISTS zitadel.users;
DROP TYPE IF EXISTS zitadel.user_type;
DROP TYPE IF EXISTS zitadel.user_state;
