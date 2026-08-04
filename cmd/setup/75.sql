ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS ios_team_id TEXT;
ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS ios_bundle_id TEXT;
ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS android_package_name TEXT;
ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS android_sha256_cert_fingerprints TEXT[];
