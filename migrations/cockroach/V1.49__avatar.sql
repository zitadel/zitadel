ALTER TABLE adminapi.iam_members ADD COLUMN avatar_key TEXT;
ALTER TABLE auth.user_grants ADD COLUMN avatar_key TEXT;
ALTER TABLE authz.user_grants ADD COLUMN avatar_key TEXT;
ALTER TABLE management.org_members ADD COLUMN avatar_key TEXT;
ALTER TABLE management.project_members ADD COLUMN avatar_key TEXT;
ALTER TABLE management.project_grant_members ADD COLUMN avatar_key TEXT;
ALTER TABLE management.user_grants ADD COLUMN avatar_key TEXT;

ALTER TABLE adminapi.iam_members ADD COLUMN user_resource_owner TEXT;
ALTER TABLE management.org_members ADD COLUMN user_resource_owner TEXT;
ALTER TABLE management.project_members ADD COLUMN user_resource_owner TEXT;
ALTER TABLE management.project_grant_members ADD COLUMN user_resource_owner TEXT;

ALTER TABLE auth.user_grants ADD COLUMN user_resource_owner TEXT;
ALTER TABLE authz.user_grants ADD COLUMN user_resource_owner TEXT;
ALTER TABLE management.user_grants ADD COLUMN user_resource_owner TEXT;