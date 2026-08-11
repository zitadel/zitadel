ALTER TABLE IF EXISTS projections.apps7_api_configs
ADD COLUMN IF NOT EXISTS minimal_introspection boolean DEFAULT false;

ALTER TABLE IF EXISTS projections.apps7_oidc_configs
ADD COLUMN IF NOT EXISTS minimal_introspection boolean DEFAULT false;
