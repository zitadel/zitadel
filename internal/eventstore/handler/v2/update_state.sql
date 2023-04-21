INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
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
    event_date = $3
    , last_updated = statement_timestamp()
;