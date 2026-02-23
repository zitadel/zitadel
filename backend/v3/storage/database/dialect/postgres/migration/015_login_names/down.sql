DROP TABLE IF EXISTS zitadel.login_names CASCADE;
DROP FUNCTION IF EXISTS zitadel.apply_domain_policy_manipulation_to_login_names();
DROP FUNCTION IF EXISTS zitadel.apply_user_insert_to_login_names();
DROP FUNCTION IF EXISTS zitadel.apply_user_update_to_login_names();
DROP FUNCTION IF EXISTS zitadel.apply_domain_manipulation_to_login_names();

DELETE FROM projections.current_states WHERE projection_name = 'zitadel.login_names';