ALTER TABLE IF EXISTS projections.idp_templates6_oauth2 ADD COLUMN IF NOT EXISTS use_pkce BOOLEAN;
ALTER TABLE IF EXISTS projections.idp_templates6_oidc ADD COLUMN IF NOT EXISTS use_pkce BOOLEAN;