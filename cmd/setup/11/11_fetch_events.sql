SELECT 
    id
    , creation_date 
FROM 
    eventstore.events 
WHERE 
    created_at IS NULL 
ORDER BY 
    event_sequence DESC
    , instance_id 
LIMIT $1 
FOR UPDATE