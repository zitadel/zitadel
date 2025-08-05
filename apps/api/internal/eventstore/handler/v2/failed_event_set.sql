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
) VALUES (
    $1
    , $2
    , $3
    , $4
    , $5
    , $6
    , $7
    , $8
    , now()
) ON CONFLICT (
    projection_name
    , aggregate_type
    , aggregate_id
    , failed_sequence
    , instance_id
) DO UPDATE SET 
    failure_count = EXCLUDED.failure_count
    , error = EXCLUDED.error
    , last_failed = EXCLUDED.last_failed
;