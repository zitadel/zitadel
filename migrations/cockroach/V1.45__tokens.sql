ALTER TABLE auth.refresh_tokens ALTER COLUMN id SET NOT NULL;
ALTER TABLE auth.refresh_tokens ALTER PRIMARY KEY USING COLUMNS(id);

ALTER TABLE auth.tokens ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE auth.tokens ALTER COLUMN application_id SET NOT NULL;
ALTER TABLE auth.tokens ALTER COLUMN user_agent_id SET NOT NULL;
ALTER TABLE auth.tokens ADD CONSTRAINT unique_token UNIQUE (user_id, application_id, user_agent_id);