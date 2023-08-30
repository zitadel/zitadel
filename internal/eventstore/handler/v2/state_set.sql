INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , aggregate_id
    , aggregate_type
    , "sequence"
    , event_date
    , position
    , last_updated
) VALUES (
    $1
    , $2
    , $3
    , $4
    , $5
    , $6
    , $7
    , statement_timestamp()
) ON CONFLICT (
    projection_name
    , instance_id
) DO UPDATE SET
    event_date = $6
    , position = $7
    , last_updated = statement_timestamp()
;