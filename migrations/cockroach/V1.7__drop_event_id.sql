ALTER TABLE eventstore.events ALTER PRIMARY KEY USING COLUMNS (event_sequence DESC);

DROP INDEX eventstore.events_id_key CASCADE;

SET sql_safe_updates = false;
ALTER TABLE eventstore.events DROP COLUMN id;
SET sql_safe_updates = true;

DROP TABLE eventstore.locks;