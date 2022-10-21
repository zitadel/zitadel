CREATE INDEX IF NOT EXISTS write_model ON eventstore.events (instance_id, aggregate_type, aggregate_id, event_type, resource_owner)
    STORING (id, aggregate_version, event_data, editor_user, editor_service);
