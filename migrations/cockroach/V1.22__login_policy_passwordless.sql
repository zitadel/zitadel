ALTER TABLE management.login_policies ADD COLUMN passwordless_type SMALLINT;
ALTER TABLE adminapi.login_policies ADD COLUMN passwordless_type SMALLINT;
ALTER TABLE auth.login_policies ADD COLUMN passwordless_type SMALLINT;
