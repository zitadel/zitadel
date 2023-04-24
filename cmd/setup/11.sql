CREATE TABLE projections.current_states (
    projection_name TEXT NOT NULL
    , instance_id TEXT NOT NULL

    , event_date TIMESTAMPTZ
    , last_updated TIMESTAMPTZ
    --aggregate is needed if multiple events are created at the same
    --point in time
    , aggregate_type TEXT
    , aggregate_id TEXT
    , event_sequence INT8

    , PRIMARY KEY (projection_name, instance_id)
);

INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , aggregate_type
    , aggregate_id
    , event_sequence
    , last_updated
) SELECT 
    cs.projection_name
    , cs.instance_id
    , e.creation_date
    , e.aggregate_type
    , e.aggregate_id
    , e.event_sequence
    , cs.timestamp
FROM 
    projections.current_sequences cs
JOIN eventstore.events e ON
    e.instance_id = cs.instance_id 
    AND e.aggregate_type = cs.aggregate_type
    AND e.event_sequence = cs.current_sequence 
    AND current_sequence = (
        SELECT 
            MAX(cs2.current_sequence)
        FROM
            projections.current_sequences cs2
        WHERE
            cs.projection_name = cs2.projection_name
            AND cs.instance_id = cs2.instance_id
    )
;