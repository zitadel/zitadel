ALTER TABLE zitadel.projections.projects ADD COLUMN project_role_assertion BOOLEAN;
ALTER TABLE zitadel.projections.projects ADD COLUMN project_role_check BOOLEAN;
ALTER TABLE zitadel.projections.projects ADD COLUMN has_project_check BOOLEAN;
ALTER TABLE zitadel.projections.projects ADD COLUMN private_labeling_setting SMALLINT;
ALTER TABLE zitadel.projections.projects ADD COLUMN sequence BIGINT;
ALTER TABLE zitadel.projections.projects RENAME COLUMN owner_id TO resource_owner;