BEGIN;

ALTER TABLE management.project_members ADD COLUMN display_name TEXT;
ALTER TABLE management.project_grant_members ADD COLUMN display_name TEXT;
ALTER TABLE management.org_members ADD COLUMN display_name TEXT;
ALTER TABLE adminapi.iam_members ADD COLUMN display_name TEXT;

COMMIT;