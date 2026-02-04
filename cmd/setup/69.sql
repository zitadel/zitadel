DO
$do$
BEGIN
   -- Convert objects to logged first (if unlogged)
   IF EXISTS (SELECT 1 FROM pg_class WHERE oid = 'cache.objects'::REGCLASS AND relpersistence = 'u') THEN
        ALTER TABLE IF EXISTS cache.objects
            SET LOGGED;
   END IF;

   -- Convert string_keys to logged (if unlogged)
   -- Drop the FK constraint first to avoid conflict with parent table conversion
   IF EXISTS (SELECT 1 FROM pg_class WHERE oid = 'cache.string_keys'::REGCLASS AND relpersistence = 'u') THEN
        ALTER TABLE IF EXISTS cache.string_keys
            DROP CONSTRAINT IF EXISTS fk_object;
        ALTER TABLE IF EXISTS cache.string_keys
            SET LOGGED;
        ALTER TABLE IF EXISTS cache.string_keys
            ADD CONSTRAINT fk_object
            FOREIGN KEY(cache_name, object_id)
            REFERENCES cache.objects(cache_name, id)
            ON DELETE CASCADE;
   END IF;

   CREATE UNLOGGED TABLE IF NOT EXISTS cache.objects_default 
        PARTITION OF cache.objects DEFAULT;

   CREATE UNLOGGED TABLE IF NOT EXISTS cache.string_keys_default 
        PARTITION OF cache.string_keys DEFAULT;
END
$do$
