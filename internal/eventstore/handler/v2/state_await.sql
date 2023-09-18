SELECT 1
FROM 
    projections.current_states
WHERE
    instance_id = $1
    AND projection_name = $2;