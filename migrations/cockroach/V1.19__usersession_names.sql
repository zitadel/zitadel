BEGIN;

ALTER TABLE auth.user_sessions ADD COLUMN user_display_name TEXT;
ALTER TABLE auth.user_sessions ADD COLUMN login_name TEXT;

COMMIT;