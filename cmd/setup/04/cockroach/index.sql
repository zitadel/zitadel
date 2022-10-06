CREATE INDEX IF NOT EXISTS write_model ON eventstore.events (instance_id, aggregate_type, aggregate_id, event_type, resource_owner)
    STORING (id, aggregate_version, previous_aggregate_sequence, creation_date, event_data, editor_user, editor_service, previous_aggregate_type_sequence);

CREATE INDEX IF NOT EXISTS active_instances ON eventstore.events (creation_date desc, instance_id) USING HASH;
