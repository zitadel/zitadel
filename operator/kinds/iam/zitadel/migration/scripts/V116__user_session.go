package scripts

const V116UserSession = `ALTER TABLE auth.user_sessions ADD COLUMN external_login_verification TIMESTAMPTZ;
ALTER TABLE auth.user_sessions ADD COLUMN selected_idp_config_id TEXT;
`
