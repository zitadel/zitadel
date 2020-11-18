ALTER TABLE management.users ADD COLUMN u2f_tokens BYTEA;
ALTER TABLE auth.users ADD COLUMN u2f_tokens JSONB;
ALTER TABLE adminapi.users ADD COLUMN u2f_tokens BYTEA;

ALTER TABLE management.users ADD COLUMN passwordless_verified_ids TEXT ARRAY;
ALTER TABLE auth.users ADD COLUMN passwordless_verified_ids TEXT ARRAY;
ALTER TABLE adminapi.users ADD COLUMN passwordless_verified_ids TEXT ARRAY;

ALTER TABLE auth.user_sessions ADD COLUMN passwordless_verification TIMESTAMPTZ;
