BEGIN;

ALTER TABLE auth.users ADD COLUMN preferred_login_name TEXT;
ALTER TABLE management.users ADD COLUMN preferred_login_name TEXT;

COMMIT;