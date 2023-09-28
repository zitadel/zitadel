INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , "position"
    , last_updated
) SELECT 
    cs.projection_name
    , cs.instance_id
    , e.created_at
    , e.position
    , cs.timestamp
FROM 
    projections.current_sequences cs
JOIN eventstore.events2 e ON
    e.instance_id = cs.instance_id 
    AND e.aggregate_type = cs.aggregate_type
    AND e.sequence = cs.current_sequence 
    AND cs.current_sequence = (
        SELECT 
            MAX(cs2.current_sequence)
        FROM
            projections.current_sequences cs2
        WHERE
            cs.projection_name = cs2.projection_name
            AND cs.instance_id = cs2.instance_id
    )
ON CONFLICT DO NOTHING;