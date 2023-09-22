INSERT INTO projections.failed_events2 (
    projection_name
    , instance_id
    , aggregate_type
    , aggregate_id
    , event_creation_date
    , failed_sequence
    , failure_count
    , error
    , last_failed
) (SELECT 
    fe.projection_name
    , fe.instance_id
    , e.aggregate_type
    , e.aggregate_id
    , e.creation_date
    , e.event_sequence
    , fe.failure_count
    , fe.error
    , fe.last_failed
FROM 
    projections.failed_events fe
JOIN eventstore.events e ON
    e.instance_id = fe.instance_id 
    AND e.event_sequence = fe.failed_sequence 
)
ON CONFLICT DO NOTHING;