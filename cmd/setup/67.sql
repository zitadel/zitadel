ALTER TABLE projections.targets2
ADD COLUMN IF NOT EXISTS payload_type smallint default 0;