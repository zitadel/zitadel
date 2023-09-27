CREATE TABLE eventstore.events2 (
    instance_id,
    aggregate_type,
    aggregate_id,
    
    event_type,
    "sequence",
    revision,
    created_at,
    payload,
    creator,
    "owner",
    
    "position",
    in_tx_order,

    PRIMARY KEY (instance_id, aggregate_type, aggregate_id, "sequence")
) AS SELECT
    instance_id,
    aggregate_type,
    aggregate_id,

    event_type,
    event_sequence,
    aggregate_version,
    creation_date,
    event_data,
    editor_user,
    resource_owner,

    COALESCE("position", creation_date::DECIMAL),
    COALESCE(in_tx_order, event_sequence)
FROM eventstore.events;