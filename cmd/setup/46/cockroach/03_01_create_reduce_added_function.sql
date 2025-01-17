CREATE OR REPLACE PROCEDURE reduce_instance_added(
    _event eventstore.events2
)
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO instances (
        id
        , name
        , creation_date
        , change_date
        , latest_position
        , latest_in_position_order
    ) VALUES (
        (_event).aggregate_id
        , (_event).payload->>'name'
        , (_event).created_at
        , (_event).created_at
        , (_event).position
        , (_event).in_tx_order
    )
    ON CONFLICT (id) DO NOTHING;
END;
$$;
