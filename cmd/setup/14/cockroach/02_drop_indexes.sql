-- not used anymore
DROP INDEX IF EXISTS eventstore.previous_sequence_unique CASCADE;
DROP INDEX IF EXISTS eventstore.events@prev_agg_type_seq_unique CASCADE;
-- dropped because equal to the PK new
DROP INDEX IF EXISTS eventstore.events@agg_type CASCADE;
DROP INDEX IF EXISTS eventstore.events@agg_type_agg_id CASCADE;
DROP INDEX IF EXISTS eventstore.events@agg_type_seq CASCADE;
DROP INDEX IF EXISTS eventstore.events@max_sequence CASCADE;
DROP INDEX IF EXISTS eventstore.events@active_instances;