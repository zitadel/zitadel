BEGIN;

ALTER TABLE auth.user_sessions ADD COLUMN user_display_name TEXT;

COMMIT;