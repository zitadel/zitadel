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
    substr(aggregate_version, 2)::SMALLINT,
    creation_date,
    event_data,
    editor_user,
    resource_owner,

    creation_date::DECIMAL,
    event_sequence
FROM eventstore.events_old;