CREATE OR REPLACE PROCEDURE reduce_instance_added("event" eventstore.events2)
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO instances (
        id
        , "name"
        , creation_date
        , change_date
        , latest_position
        , latest_in_position_order
    ) VALUES (
        event.aggregate_id
        , event.payload->>'name'
        , event.created_at
        , event.created_at
        , event.position
        , event.in_tx_order::INT2
    )
    ON CONFLICT (id) DO NOTHING;
END;
$$;
