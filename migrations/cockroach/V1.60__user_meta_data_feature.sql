ALTER TABLE adminapi.features ADD COLUMN metadata_user BOOLEAN;
ALTER TABLE auth.features ADD COLUMN metadata_user BOOLEAN;
ALTER TABLE authz.features ADD COLUMN metadata_user BOOLEAN;
ALTER TABLE management.features ADD COLUMN metadata_user BOOLEAN;