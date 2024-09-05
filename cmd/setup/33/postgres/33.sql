CREATE INDEX CONCURRENTLY IF NOT EXISTS es_active_instances_idx  ON eventstore.events2 (aggregate_type, event_type) WHERE aggregate_type = 'instance' AND event_type IN ('instance.added', 'instance.removed');
CREATE INDEX CONCURRENTLY IF NOT EXISTS es_current_sequence_idx ON eventstore.events2 ("sequence" DESC, aggregate_id, aggregate_type, instance_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS es_inst_agg_typ_id_idx ON eventstore.events2 (instance_id, aggregate_id, aggregate_type, "position");
CREATE INDEX CONCURRENTLY IF NOT EXISTS es_inst_agg_typ_event_typ_idx ON eventstore.events2 (instance_id, aggregate_type, event_type, "position");

-- TODO: remove old indexes