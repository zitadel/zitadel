ALTER TABLE projections.targets2
ADD COLUMN IF NOT EXISTS payload_type smallint default 0;
ALTER TABLE projections.authn_keys2
ADD COLUMN IF NOT EXISTS fingerprint text;