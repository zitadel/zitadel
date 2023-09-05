BEGIN;
ALTER TABLE eventstore.events ALTER COLUMN editor_service DROP NOT NULL;

-- not used anymore
DROP INDEX IF EXISTS eventstore.previous_sequence_unique CASCADE;
DROP INDEX IF EXISTS eventstore.prev_agg_type_seq_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events_event_sequence_instance_id_key CASCADE;
-- dropped because equal to the PK new
DROP INDEX IF EXISTS eventstore.agg_type CASCADE;
DROP INDEX IF EXISTS eventstore.agg_type_agg_id CASCADE;
DROP INDEX IF EXISTS eventstore.agg_type_seq CASCADE;
DROP INDEX IF EXISTS eventstore.max_sequence CASCADE;
DROP INDEX IF EXISTS eventstore.active_instances;
COMMIT;

BEGIN;
ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS "position" BIGINT NOT NULL;
COMMIT;

BEGIN;
ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS in_tx_order SMALLINT NOT NULL;
COMMIT;

BEGIN;
ALTER TABLE eventstore.events DROP CONSTRAINT IF EXISTS events_pkey;
ALTER TABLE eventstore.events ADD PRIMARY KEY (instance_id, aggregate_type, aggregate_id, event_sequence);
COMMIT;

BEGIN;
CREATE INDEX IF NOT EXISTS es_handler_idx ON eventstore.events (instance_id, aggregate_type, event_type);
CREATE INDEX IF NOT EXISTS es_agg_id_event_idx ON eventstore.events (aggregate_type, aggregate_id, event_type);
CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events (created_at DESC, instance_id);
CREATE INDEX IF NOT EXISTS es_global ON eventstore.events (aggregate_type, aggregate_id);
CREATE INDEX IF NOT EXISTS es_inst_agg_type_agg_id_cd ON eventstore.events (instance_id, aggregate_type, aggregate_id, event_sequence);
COMMIT;

