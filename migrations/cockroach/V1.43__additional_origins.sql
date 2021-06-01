ALTER TABLE auth.applications ADD COLUMN additional_origins TEXT ARRAY;
ALTER TABLE authz.applications ADD COLUMN additional_origins TEXT ARRAY;
ALTER TABLE management.applications ADD COLUMN additional_origins TEXT ARRAY;