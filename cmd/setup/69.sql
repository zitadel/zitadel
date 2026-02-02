ALTER TABLE IF EXISTS cache.objects
    SET LOGGED;

ALTER TABLE IF EXISTS cache.string_keys
    SET LOGGED;

CREATE UNLOGGED TABLE IF NOT EXISTS cache.objects_default 
    PARTITION OF cache.objects DEFAULT;

CREATE UNLOGGED TABLE IF NOT EXISTS cache.string_keys_default 
    PARTITION OF cache.string_keys DEFAULT;