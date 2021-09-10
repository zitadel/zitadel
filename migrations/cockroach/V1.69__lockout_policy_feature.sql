ALTER TABLE adminapi.features ADD COLUMN lockout_policy BOOLEAN;
ALTER TABLE auth.features ADD COLUMN lockout_policy BOOLEAN;
ALTER TABLE authz.features ADD COLUMN lockout_policy BOOLEAN;
ALTER TABLE management.features ADD COLUMN lockout_policy BOOLEAN;