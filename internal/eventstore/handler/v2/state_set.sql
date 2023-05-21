INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , aggregate_type
    , aggregate_id
    , event_sequence
    , last_updated
) VALUES (
    $1
    , $2
    , $3
    , $4
    , $5
    , $6
    , statement_timestamp()
) ON CONFLICT (
    projection_name
    , instance_id
) DO UPDATE SET
    event_date = $3
    , aggregate_type = $4
    , aggregate_id = $5
    , event_sequence = $6
    , last_updated = statement_timestamp()
;