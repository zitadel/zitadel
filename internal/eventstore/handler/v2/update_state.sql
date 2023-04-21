INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_timestamp
    , last_updated
) VALUES (
    $1
    , $2
    , $3
    , statement_time()
) ON CONFLICT (
    projection_name
    , instance_id
) DO UPDATE SET
    event_timestamp = EXCLUEDED.event_timestamp
    , last_updated = EXCLUEDED.last_updated
;