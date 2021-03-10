ALTER TABLE management.login_policies ADD COLUMN default_policy BOOLEAN;
ALTER TABLE adminapi.login_policies ADD COLUMN default_policy BOOLEAN;
ALTER TABLE auth.login_policies ADD COLUMN default_policy BOOLEAN;
