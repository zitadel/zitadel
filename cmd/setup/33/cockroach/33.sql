CREATE INDEX CONCURRENTLY IF NOT EXISTS es_active_instances_events  ON eventstore.events2 (aggregate_type, event_type) WHERE aggregate_type = 'instance' AND event_type IN ('instance.added', 'instance.removed'); 
CREATE INDEX CONCURRENTLY IF NOT EXISTS es_current_sequence ON eventstore.events2 ("sequence" DESC, aggregate_id, aggregate_type, instance_id);

CREATE INDEX es_inst_agg_typ_event_typ_pos_idx ON eventstore.events2 (instance_id, aggregate_type, event_type, "position") STORING (revision, created_at, payload, creator, owner, in_tx_order);
CREATE INDEX es_inst_agg_typ_pos_idx ON eventstore.events2 (instance_id, aggregate_type, "position") STORING (event_type, revision, created_at, payload, creator, owner, in_tx_order);

-- TODO: remove old indexes