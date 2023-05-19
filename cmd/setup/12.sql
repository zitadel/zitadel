BEGIN;
ALTER TABLE eventstore.events ALTER PRIMARY KEY USING COLUMNS (instance_id, aggregate_type, aggregate_id, event_sequence DESC);
COMMIT;

BEGIN;
ALTER TABLE eventstore.events ALTER COLUMN editor_service DROP NOT NULL;

-- not used anymore
DROP INDEX IF EXISTS eventstore.events@previous_sequence_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events@prev_agg_type_seq_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events@events_event_sequence_instance_id_key CASCADE;
-- dropped because equal to the PK new
DROP INDEX IF EXISTS eventstore.events@max_sequence CASCADE;
COMMIT;
