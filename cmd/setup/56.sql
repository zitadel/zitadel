ALTER TABLE IF EXISTS projections.idp_templates6_saml ADD COLUMN IF NOT EXISTS federated_logout_enabled BOOLEAN DEFAULT FALSE;
