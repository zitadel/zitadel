ALTER TABLE IF EXISTS projections.idp_templates6_saml ADD COLUMN IF NOT EXISTS name_id_format SMALLINT;
ALTER TABLE IF EXISTS projections.idp_templates6_saml ADD COLUMN IF NOT EXISTS transient_mapping_attribute_name TEXT;
