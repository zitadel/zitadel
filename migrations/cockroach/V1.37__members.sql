ALTER TABLE management.org_members ADD COLUMN preferred_login_name TEXT;
ALTER TABLE management.project_members ADD COLUMN preferred_login_name TEXT;
ALTER TABLE management.project_grant_members ADD COLUMN preferred_login_name TEXT;

ALTER TABLE adminapi.iam_members ADD COLUMN preferred_login_name TEXT;
