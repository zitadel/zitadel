BEGIN;

UPDATE auth.authn_keys SET object_type = 1 where object_type = 0;
UPDATE management.authn_keys SET object_type = 1 where object_type = 0;

COMMIT;
