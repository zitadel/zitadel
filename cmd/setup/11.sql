ALTER TABLE eventstore.events ADD COLUMN created_at TIMESTAMPTZ;

UPDATE eventstore.events SET created_at = creation_date WHERE created_at IS NULL;

ALTER TABLE eventstore.events ALTER COLUMN created_at SET NOT NULL DEFAULT clock_timestamp();