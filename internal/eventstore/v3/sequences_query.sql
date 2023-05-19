SELECT
    instance_id
    , aggregate_type
    , aggregate_id
    , event_sequence
FROM 
    eventstore.events
WHERE (
    instance_id
    , aggregate_type
    , aggregate_id
    , event_sequence
    )  IN (
        SELECT (
            instance_id
            , aggregate_type
            , aggregate_id
            , MAX(event_sequence)
        ) FROM 
            eventstore.events 
        WHERE 
            %s 
        GROUP BY 
            instance_id
            , aggregate_type
            , aggregate_id
    )
FOR UPDATE;