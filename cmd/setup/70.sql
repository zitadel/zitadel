ALTER TABLE IF EXISTS system.assets
    DROP COLUMN IF EXISTS hash;

ALTER TABLE IF EXISTS system.assets
    ADD COLUMN hash TEXT GENERATED ALWAYS AS (encode(sha256(data), 'hex')) STORED;
