BEGIN;
ALTER TABLE eventstore.events ALTER COLUMN editor_service DROP NOT NULL;

-- not used anymore
DROP INDEX IF EXISTS eventstore.events@previous_sequence_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events@prev_agg_type_seq_unique CASCADE;
-- dropped because equal to the PK new
DROP INDEX IF EXISTS eventstore.events@agg_type CASCADE;
DROP INDEX IF EXISTS eventstore.events@agg_type_agg_id CASCADE;
DROP INDEX IF EXISTS eventstore.events@agg_type_seq CASCADE;
DROP INDEX IF EXISTS eventstore.events@max_sequence CASCADE;
DROP INDEX IF EXISTS eventstore.events@active_instances;
COMMIT;

-- BEGIN;
-- ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS "position" BIGINT NOT NULL;
-- COMMIT;

-- BEGIN;
-- ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS commit_order DECIMAL NOT NULL;
-- COMMIT;

-- TOOD: create a new migration step which removes created_at, previous_aggregate_*

BEGIN;
ALTER TABLE eventstore.events ALTER PRIMARY KEY USING COLUMNS (instance_id, aggregate_type, aggregate_id, event_sequence DESC);
COMMIT;

BEGIN;
DROP INDEX IF EXISTS eventstore.events@events_event_sequence_instance_id_key CASCADE;
COMMIT;

BEGIN;
CREATE INDEX IF NOT EXISTS es_handler_idx ON eventstore.events (instance_id, "position", aggregate_type, event_type) INCLUDE (created_at, event_data, editor_user, resource_owner, aggregate_version, in_tx_order);
CREATE INDEX IF NOT EXISTS es_agg_id_event_idx ON eventstore.events (aggregate_type, aggregate_id, event_type) INCLUDE (created_at, event_data, editor_user, resource_owner, aggregate_version, "position", in_tx_order);
CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events (created_at, instance_id) USING HASH;
CREATE INDEX IF NOT EXISTS es_global ON eventstore.events (aggregate_type, aggregate_id) INCLUDE (created_at, event_type, event_data, editor_user, resource_owner, aggregate_version, "position", in_tx_order);
CREATE INDEX IF NOT EXISTS es_inst_agg_type_agg_id_cd ON eventstore.events (instance_id, aggregate_type, aggregate_id, event_sequence) INCLUDE (created_at, event_type, event_data, editor_user, resource_owner, aggregate_version, "position", in_tx_order);
COMMIT;
