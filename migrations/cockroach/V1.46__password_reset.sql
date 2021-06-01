ALTER TABLE adminapi.features ADD COLUMN login_policy_password_reset BOOLEAN;
ALTER TABLE auth.features ADD COLUMN login_policy_password_reset BOOLEAN;
ALTER TABLE authz.features ADD COLUMN login_policy_password_reset BOOLEAN;
ALTER TABLE management.features ADD COLUMN login_policy_password_reset BOOLEAN;


ALTER TABLE auth.login_policies ADD COLUMN hide_password_reset BOOLEAN;
ALTER TABLE adminapi.login_policies ADD COLUMN hide_password_reset BOOLEAN;
ALTER TABLE management.login_policies ADD COLUMN hide_password_reset BOOLEAN;
