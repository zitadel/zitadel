package scripts

const V18UsernameChange = `
BEGIN;

ALTER TABLE management.users ADD COLUMN username_change_required BOOLEAN;
ALTER TABLE auth.users ADD COLUMN username_change_required BOOLEAN;
ALTER TABLE adminapi.users ADD COLUMN username_change_required BOOLEAN;

COMMIT;
`
