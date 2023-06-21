SELECT 
    event_date
    , aggregate_type
    , aggregate_id
    , event_sequence
FROM 
    projections.current_states
WHERE
    instance_id = $1
    AND projection_name = $2
FOR UPDATE NOWAIT