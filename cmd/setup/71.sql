ALTER TABLE IF EXISTS projections.idp_templates6_jwt
ADD COLUMN IF NOT EXISTS audience text;
