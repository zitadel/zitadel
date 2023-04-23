INSERT INTO projections.failed_events (
    projection_name
    , failed_sequence
    , failure_count
    , error
    , instance_id
    , last_failed
) VALUES ($1, $2, $3, $4, $5, statement_timestamp()) 
ON CONFLICT (
    projection_name
    , failed_sequence
    , instance_id
) DO UPDATE SET 
    failure_count = EXCLUDED.failure_count
    , error = EXCLUDED.error
    , last_failed = EXCLUDED.last_failed
;