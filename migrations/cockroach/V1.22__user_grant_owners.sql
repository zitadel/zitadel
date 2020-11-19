ALTER TABLE management.user_grants ADD COLUMN project_owner STRING;
ALTER TABLE management.user_grants RENAME COLUMN resource_owner TO grant_owner;

ALTER TABLE auth.user_grants ADD COLUMN project_owner STRING;
ALTER TABLE auth.user_grants RENAME COLUMN resource_owner TO grant_owner;

ALTER TABLE authz.user_grants ADD COLUMN project_owner STRING;
ALTER TABLE authz.user_grants RENAME COLUMN resource_owner TO grant_owner;
