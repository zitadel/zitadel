ALTER TABLE adminapi.users ADD COLUMN passwordless_init_required boolean;
ALTER TABLE auth.users ADD COLUMN passwordless_init_required boolean;
ALTER TABLE management.users ADD COLUMN passwordless_init_required boolean;

BEGIN;
ALTER TABLE adminapi.users RENAME COLUMN password_set TO password_init_required;
UPDATE adminapi.users set password_init_required = NOT password_init_required;

ALTER TABLE auth.users RENAME COLUMN password_set TO password_init_required;
UPDATE auth.users set password_init_required = NOT password_init_required;

ALTER TABLE management.users RENAME COLUMN password_set TO password_init_required;
UPDATE management.users set password_init_required = NOT password_init_required;

COMMIT;