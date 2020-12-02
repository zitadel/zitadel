ALTER TABLE management.login_policies ADD COLUMN passwordless_type SMALLINT;
ALTER TABLE adminapi.login_policies ADD COLUMN passwordless_type SMALLINT;
ALTER TABLE auth.login_policies ADD COLUMN passwordless_type SMALLINT;

ALTER TABLE management.users ADD COLUMN u2f_tokens BYTEA;
ALTER TABLE auth.users ADD COLUMN u2f_tokens BYTEA;
ALTER TABLE adminapi.users ADD COLUMN u2f_tokens BYTEA;

ALTER TABLE management.users ADD COLUMN passwordless_tokens BYTEA;
ALTER TABLE auth.users ADD COLUMN passwordless_tokens BYTEA;
ALTER TABLE adminapi.users ADD COLUMN passwordless_tokens BYTEA;

ALTER TABLE auth.user_sessions ADD COLUMN passwordless_verification TIMESTAMPTZ;
