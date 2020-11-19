ALTER TABLE management.user_grants ADD COLUMN project_owner STRING;

ALTER TABLE auth.user_grants ADD COLUMN project_owner STRING;

ALTER TABLE authz.user_grants ADD COLUMN project_owner STRING;
