ALTER TABLE zitadel.projections.projects ADD COLUMN project_role_assertion BOOLEAN;
ALTER TABLE zitadel.projections.projects ADD COLUMN project_role_check BOOLEAN;
ALTER TABLE zitadel.projections.projects ADD COLUMN has_project_check BOOLEAN;
ALTER TABLE zitadel.projections.projects ADD COLUMN private_labeling_setting SMALLINT;
ALTER TABLE zitadel.projections.projects ADD COLUMN sequence BIGINT;
ALTER TABLE zitadel.projections.projects RENAME COLUMN owner_id TO resource_owner;


CREATE TABLE zitadel.projections.project_grants (
      project_id TEXT,
      grant_id TEXT,
      creation_date TIMESTAMPTZ,
      change_date TIMESTAMPTZ,
      resource_owner TEXT,
      state INT2,
      sequence BIGINT,

      project_name TEXT,
      granted_org_id TEXT,
      granted_org_name TEXT,
      granted_role_keys TEXT Array,
      resource_owner_name TEXT,
      creator_id TEXT,

      PRIMARY KEY (grant_id)
      LEFT JOIN projections.orgs o ON o.id = u.org_id
);
