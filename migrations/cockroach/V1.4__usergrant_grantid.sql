BEGIN;

ALTER TABLE management.user_grants ADD COLUMN grant_id TEXT;
ALTER TABLE auth.user_grants ADD COLUMN grant_id TEXT;
ALTER TABLE authz.user_grants ADD COLUMN grant_id TEXT;

COMMIT;