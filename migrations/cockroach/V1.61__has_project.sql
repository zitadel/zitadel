ALTER TABLE management.projects ADD COLUMN has_project_check BOOLEAN;

ALTER TABLE authz.applications ADD COLUMN has_project_check BOOLEAN;
ALTER TABLE auth.applications ADD COLUMN has_project_check BOOLEAN;
ALTER TABLE management.applications ADD COLUMN has_project_check BOOLEAN;

CREATE TABLE auth.org_project_mapping (
     org_id TEXT,
     project_id TEXT,

     PRIMARY KEY (org_id, project_id)
);