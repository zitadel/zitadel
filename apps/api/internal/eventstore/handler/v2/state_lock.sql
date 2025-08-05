INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , last_updated
) VALUES (
    $1
    , $2
    , now()
) ON CONFLICT DO NOTHING;