BEGIN;


ALTER TABLE auth.users ADD COLUMN login_names TEXT ARRAY;
ALTER TABLE management.users ADD COLUMN login_names TEXT ARRAY;

COMMIT;