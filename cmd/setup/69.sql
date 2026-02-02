DO
$do$
BEGIN
   IF EXISTS (SELECT 1 FROM pg_class WHERE oid = 'cache.objects'::REGCLASS::OID AND relpersistence = 'u') THEN
        ALTER TABLE IF EXISTS cache.objects
            SET LOGGED;
   END IF;
   IF EXISTS (SELECT 1 FROM pg_class WHERE oid = 'cache.string_keys'::REGCLASS::OID AND relpersistence = 'u') THEN
        ALTER TABLE IF EXISTS cache.string_keys
            SET LOGGED;
   END IF;

   CREATE UNLOGGED TABLE IF NOT EXISTS cache.objects_default 
        PARTITION OF cache.objects DEFAULT;

   CREATE UNLOGGED TABLE IF NOT EXISTS cache.string_keys_default 
        PARTITION OF cache.string_keys DEFAULT;
END
$do$
