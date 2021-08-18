ALTER TABLE management.projects ADD COLUMN private_labeling_setting SMALLINT;

ALTER TABLE authz.applications ADD COLUMN private_labeling_setting SMALLINT;
ALTER TABLE auth.applications ADD COLUMN private_labeling_setting SMALLINT;
ALTER TABLE management.applications ADD COLUMN private_labeling_setting SMALLINT;
