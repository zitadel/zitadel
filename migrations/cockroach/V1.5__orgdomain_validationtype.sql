BEGIN;

ALTER TABLE management.org_domains ADD COLUMN validation_type SMALLINT;

COMMIT;