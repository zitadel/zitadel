WITH existing AS (
    %s
) SELECT 
    e.instance_id
    , e.owner
    , e.aggregate_type
    , e.aggregate_id
    , e.sequence
FROM 
    eventstore.events2 e
JOIN 
    existing 
ON 
    e.instance_id = existing.instance_id
    AND e.aggregate_type = existing.aggregate_type
    AND e.aggregate_id = existing.aggregate_id
    AND e.sequence = existing.sequence
FOR NO KEY UPDATE;