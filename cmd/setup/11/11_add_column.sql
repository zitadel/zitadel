BEGIN;
-- create table with empty created_at
ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT NULL;
-- set column rules
ALTER TABLE eventstore.events ALTER COLUMN created_at SET DEFAULT clock_timestamp();
COMMIT;