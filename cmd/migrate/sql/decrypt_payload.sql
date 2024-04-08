INSERT INTO eventstore.decrypted_payload (
    instance_id,
    aggregate_type,
    aggregate_id,
    "sequence",

    decrypted_payload,
) SELECT
    instance_id,
    aggregate_type,
    aggregate_id,
    "sequence",
    -- TODO: switch
FROM 
    eventstore.events2 
WHERE
    -- TODO: instance id
    AND position < $2
    AND (
        (aggregate_type = $1 AND event_type = ANY($2))
        OR (aggregate_type = $1 AND event_type = ANY($2))
        OR (aggregate_type = $1 AND event_type = ANY($2))
        OR (aggregate_type = $1 AND event_type = ANY($2))
    )