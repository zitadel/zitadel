ALTER TABLE management.projects ADD COLUMN private_labeling_setting SMALLINT;

ALTER TABLE authz.applications ADD COLUMN private_labeling_setting BOOLEAN;
ALTER TABLE auth.applications ADD COLUMN private_labeling_setting BOOLEAN;
ALTER TABLE management.applications ADD COLUMN private_labeling_setting BOOLEAN;

ALTER TABLE authz.applications ADD COLUMN resource_owner STRING;
ALTER TABLE auth.applications ADD COLUMN resource_owner STRING;
ALTER TABLE management.applications ADD COLUMN resource_owner STRING;