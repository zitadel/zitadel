ALTER TABLE IF EXISTS projections.apps7_oidc_configs ADD COLUMN IF NOT EXISTS allowed_scope_prefixes TEXT[];
