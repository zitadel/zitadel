ALTER TABLE zitadel.users
    ADD COLUMN IF NOT EXISTS password_locked_at TIMESTAMPTZ;
