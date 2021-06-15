ALTER TABLE adminapi.idp_configs ADD COLUMN auth_connector_base_url TEXT;
ALTER TABLE adminapi.idp_configs ADD COLUMN auth_connector_backend_connector_id TEXT;
ALTER TABLE auth.idp_configs ADD COLUMN auth_connector_base_url TEXT;
ALTER TABLE auth.idp_configs ADD COLUMN auth_connector_backend_connector_id TEXT;
ALTER TABLE management.idp_configs ADD COLUMN auth_connector_base_url TEXT;
ALTER TABLE management.idp_configs ADD COLUMN auth_connector_backend_connector_id TEXT;