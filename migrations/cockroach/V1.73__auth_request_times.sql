ALTER TABLE auth.auth_requests ADD COLUMN creation_data TIMESTAMPTZ;
ALTER TABLE auth.auth_requests ADD COLUMN change_data TIMESTAMPTZ;
