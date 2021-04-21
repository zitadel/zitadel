ALTER TABLE adminapi.features ADD COLUMN custom_domain BOOLEAN;
ALTER TABLE auth.features ADD COLUMN custom_domain BOOLEAN;
ALTER TABLE authz.features ADD COLUMN custom_domain BOOLEAN;
ALTER TABLE management.features ADD COLUMN custom_domain BOOLEAN;