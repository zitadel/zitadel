ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS login_version SMALLINT;
ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS login_base_uri TEXT;
