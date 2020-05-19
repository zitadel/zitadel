BEGIN;

ALTER TABLE management.users
    ADD COLUMN password_set BOOLEAN,
    ADD COLUMN password_change_required BOOLEAN,
    ADD COLUMN mfa_max_set_up SMALLINT,
    ADD COLUMN mfa_init_skipped TIMESTAMPTZ;

COMMIT;