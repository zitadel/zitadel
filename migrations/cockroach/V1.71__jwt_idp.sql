ALTER TABLE auth.idp_configs ADD COLUMN jwt_endpoint TEXT;
ALTER TABLE auth.idp_configs ADD COLUMN jwt_keys_endpoint TEXT;
ALTER TABLE auth.idp_configs ADD COLUMN jwt_header_name TEXT;

ALTER TABLE adminapi.idp_configs ADD COLUMN jwt_endpoint TEXT;
ALTER TABLE adminapi.idp_configs ADD COLUMN jwt_keys_endpoint TEXT;
ALTER TABLE adminapi.idp_configs ADD COLUMN jwt_header_name TEXT;

ALTER TABLE management.idp_configs ADD COLUMN jwt_endpoint TEXT;
ALTER TABLE management.idp_configs ADD COLUMN jwt_keys_endpoint TEXT;
ALTER TABLE management.idp_configs ADD COLUMN jwt_header_name TEXT;
