BEGIN;

ALTER TABLE management.project_grants ADD COLUMN resource_owner_name TEXT;

COMMIT;