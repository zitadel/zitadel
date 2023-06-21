with existing as (
    SELECT
        instance_id
        , aggregate_type
        , aggregate_id
        , MAX(event_sequence) event_sequence
    FROM 
        eventstore.events existing
    WHERE 
        %s
    GROUP BY 
        instance_id
        , aggregate_type
        , aggregate_id
) SELECT 
    e.instance_id
    , e.resource_owner
    , e.aggregate_type
    , e.aggregate_id
    , e.event_sequence
FROM 
    eventstore.events e
JOIN 
    existing 
ON 
        e.instance_id = existing.instance_id
        AND e.aggregate_type = existing.aggregate_type
        AND e.aggregate_id = existing.aggregate_id
        AND e.event_sequence = existing.event_sequence
FOR UPDATE;