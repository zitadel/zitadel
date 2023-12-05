INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , aggregate_id
    , aggregate_type
    , "sequence"
    , event_date
    , "position"
    , last_updated
    , filter_offset
) VALUES (
    $1
    , $2
    , $3
    , $4
    , $5
    , $6
    , $7
    , now()
    , $8
) ON CONFLICT (
    projection_name
    , instance_id
) DO UPDATE SET
    aggregate_id = $3
    , aggregate_type = $4
    , "sequence" = $5
    , event_date = $6
    , "position" = $7
    , last_updated = statement_timestamp()
    , filter_offset = $8
;