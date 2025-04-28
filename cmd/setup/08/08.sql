CREATE INDEX IF NOT EXISTS inst_refresh_tkn_idx ON auth.tokens(instance_id, refresh_token_id);
CREATE INDEX IF NOT EXISTS inst_app_tkn_idx ON auth.tokens(instance_id, application_id);
CREATE INDEX IF NOT EXISTS inst_ro_tkn_idx ON auth.tokens(instance_id, resource_owner);
DROP INDEX IF EXISTS auth.user_user_agent_idx;
CREATE INDEX IF NOT EXISTS inst_usr_agnt_tkn_idx ON auth.tokens(instance_id, user_id, user_agent_id);