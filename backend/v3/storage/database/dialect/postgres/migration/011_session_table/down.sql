DROP TABLE IF EXISTS zitadel.session_metadata CASCADE;
DROP TABLE IF EXISTS zitadel.session_user_agents CASCADE;
DROP TYPE IF EXISTS zitadel.session_factor_type CASCADE;
DROP TABLE IF EXISTS zitadel.session_factors CASCADE;
DROP TABLE IF EXISTS zitadel.sessions CASCADE;
DROP FUNCTION IF EXISTS zitadel.update_expiration();