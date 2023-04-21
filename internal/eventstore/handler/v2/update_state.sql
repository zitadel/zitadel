INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_timestamp
    , last_updated
) VALUES (
    $1
    , $2
    , $3
    , statement_timestamp()
) ON CONFLICT (
    projection_name
    , instance_id
) DO UPDATE SET
    event_timestamp = $3
    , last_updated = statement_timestamp()
;