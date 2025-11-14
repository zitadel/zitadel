ALTER TABLE IF EXISTS auth.user_sessions ADD COLUMN IF NOT EXISTS id TEXT;

CREATE INDEX IF NOT EXISTS user_session_id ON auth.user_sessions (id, instance_id);