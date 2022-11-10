ALTER TABLE auth.users ADD COLUMN owner_removed BOOL DEFAULT false;
ALTER TABLE auth.user_external_idps ADD COLUMN owner_removed BOOL DEFAULT false;
ALTER TABLE auth.org_project_mapping ADD COLUMN owner_removed BOOL DEFAULT false;
ALTER TABLE auth.idp_configs ADD COLUMN owner_removed BOOL DEFAULT false;
ALTER TABLE auth.idp_providers ADD COLUMN owner_removed BOOL DEFAULT false;

ALTER TABLE adminapi.styling ADD COLUMN owner_removed BOOL DEFAULT false;
