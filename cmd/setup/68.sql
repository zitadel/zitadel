ALTER TABLE IF EXISTS projections.targets2
ADD COLUMN IF NOT EXISTS payload_type smallint default 0;

ALTER TABLE IF EXISTS projections.authn_keys2
ADD COLUMN IF NOT EXISTS fingerprint text;

ALTER TABLE IF EXISTS projections.authn_keys2
ALTER COLUMN expiration DROP NOT NULL;
