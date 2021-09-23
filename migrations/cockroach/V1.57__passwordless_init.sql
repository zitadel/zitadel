ALTER TABLE adminapi.users ADD COLUMN passwordless_init_required boolean;
ALTER TABLE auth.users ADD COLUMN passwordless_init_required boolean;
ALTER TABLE management.users ADD COLUMN passwordless_init_required boolean;

ALTER TABLE adminapi.users ADD COLUMN password_init_required BOOLEAN;
UPDATE adminapi.users set password_init_required = NOT password_set;

ALTER TABLE auth.users ADD COLUMN password_init_required BOOLEAN;
UPDATE auth.users set password_init_required = NOT password_set;

ALTER TABLE management.users ADD COLUMN password_init_required BOOLEAN;
UPDATE management.users set password_init_required = NOT password_set;
