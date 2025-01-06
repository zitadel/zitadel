CREATE OR REPLACE FUNCTION reduce_instance_added("event" RECORD)
RETURNS VOID
LANGUAGE SQL
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
        (event).aggregate_id
        , (event).payload->>'name'
        , (event).created_at
        , (event).created_at
        , (event).position
        , (event).in_tx_order::INT2
    )
    ON CONFLICT (id) DO NOTHING;
END;
$$;
