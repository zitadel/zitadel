ALTER TABLE IF EXISTS projections.idp_templates6_oidc ADD COLUMN IF NOT EXISTS token_endpoint_auth_method SMALLINT DEFAULT 0;
ALTER TABLE IF EXISTS projections.idp_templates6_oidc ADD COLUMN IF NOT EXISTS private_key JSONB;
ALTER TABLE IF EXISTS projections.idp_templates6_oidc ADD COLUMN IF NOT EXISTS key_id TEXT;
-- client secret is not set for the private_key_jwt auth method
ALTER TABLE IF EXISTS projections.idp_templates6_oidc ALTER COLUMN client_secret DROP NOT NULL;
