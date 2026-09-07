ALTER TABLE IF EXISTS projections.idp_templates6_oauth2 ADD COLUMN IF NOT EXISTS authorization_parameters JSONB;
ALTER TABLE IF EXISTS projections.idp_templates6_oauth2 ADD COLUMN IF NOT EXISTS forwarded_parameters TEXT[];
ALTER TABLE IF EXISTS projections.idp_templates6_oidc ADD COLUMN IF NOT EXISTS authorization_parameters JSONB;
ALTER TABLE IF EXISTS projections.idp_templates6_oidc ADD COLUMN IF NOT EXISTS forwarded_parameters TEXT[];
