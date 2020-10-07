CREATE TABLE auth.project_roles (
      project_id TEXT,
      role_key TEXT,
      display_name TEXT,
      resource_owner TEXT,
      org_id TEXT,
      group_name TEXT,

      creation_date TIMESTAMPTZ,
      sequence BIGINT,

      PRIMARY KEY (org_id, project_id, role_key)
);

ALTER TABLE authz.user_grants ADD COLUMN org_primary_domain TEXT;
ALTER TABLE auth.user_grants ADD COLUMN org_primary_domain TEXT;
ALTER TABLE management.user_grants ADD COLUMN org_primary_domain TEXT;