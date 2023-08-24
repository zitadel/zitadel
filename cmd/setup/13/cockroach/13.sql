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
-- ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS position DECIMAL NOT NULL;
-- COMMIT;

BEGIN;
ALTER TABLE eventstore.events ALTER PRIMARY KEY USING COLUMNS (instance_id, aggregate_type, aggregate_id, event_sequence DESC);
COMMIT;

BEGIN;
DROP INDEX IF EXISTS eventstore.events@events_event_sequence_instance_id_key CASCADE;
COMMIT;

BEGIN;
CREATE INDEX IF NOT EXISTS es_handler_idx ON eventstore.events (instance_id, aggregate_type, event_type) INCLUDE (aggregate_version, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_user, editor_service, resource_owner);
CREATE INDEX IF NOT EXISTS es_agg_id_event_idx ON eventstore.events (aggregate_type, aggregate_id, event_type) INCLUDE (aggregate_version, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_user, editor_service, resource_owner);
CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events (aggregate_type, event_type);
CREATE INDEX IF NOT EXISTS es_global ON eventstore.events (aggregate_type, aggregate_id) INCLUDE (event_type, aggregate_version, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_user, editor_service, resource_owner);
CREATE INDEX IF NOT EXISTS es_inst_agg_type_agg_id_cd ON eventstore.events (instance_id, aggregate_type, aggregate_id) INCLUDE (event_type, aggregate_version, previous_aggregate_sequence, previous_aggregate_type_sequence, event_data, editor_user, editor_service, resource_owner);
COMMIT;
