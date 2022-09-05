CREATE INDEX IF NOT EXISTS write_model ON eventstore.events (instance_id, aggregate_type, aggregate_id, event_type, resource_owner);

CREATE INDEX IF NOT EXISTS active_instances ON eventstore.events (creation_date, instance_id);
