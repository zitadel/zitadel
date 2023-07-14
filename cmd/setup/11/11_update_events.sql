UPDATE eventstore.events SET 
    created_at =  creation_date 
FROM (
    SELECT 
        e.event_sequence as seq
        , e.instance_id as i_id
        , e.creation_date as cd
    FROM 
        eventstore.events e
    WHERE 
        created_at IS NULL
    ORDER BY 
        event_sequence ASC
        , instance_id 
    LIMIT $1
) AS e 
WHERE 
    e.seq = eventstore.events.event_sequence 
    AND e.i_id = eventstore.events.instance_id 
    AND e.cd = eventstore.events.creation_date
;