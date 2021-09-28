ALTER TABLE auth.auth_requests ADD COLUMN creation_date TIMESTAMPTZ;
ALTER TABLE auth.auth_requests ADD COLUMN change_date TIMESTAMPTZ;
