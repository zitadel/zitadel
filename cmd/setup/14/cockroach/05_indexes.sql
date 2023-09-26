DROP INDEX IF EXISTS eventstore.events@write_model;
CREATE INDEX IF NOT EXISTS es_handler_idx ON eventstore.events (aggregate_type, event_type, instance_id);
CREATE INDEX IF NOT EXISTS es_agg_wm ON eventstore.events (aggregate_type, aggregate_id, event_type, "position");
CREATE INDEX IF NOT EXISTS es_wm ON eventstore.events (instance_id, resource_owner, aggregate_type, aggregate_id, event_type);
CREATE INDEX IF NOT EXISTS es_active_instances ON eventstore.events (creation_date DESC, instance_id) USING HASH;
