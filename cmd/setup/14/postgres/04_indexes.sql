CREATE INDEX IF NOT EXISTS es_handler_idx ON eventstore.events2 (instance_id, aggregate_type, event_type);
CREATE INDEX IF NOT EXISTS es_agg_wm ON eventstore.events2 (instance_id, aggregate_type, aggregate_id, event_type, "position");
CREATE INDEX IF NOT EXISTS es_wm ON eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, event_type);
CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events2 (created_at DESC, instance_id) USING HASH;