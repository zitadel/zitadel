CREATE OR REPLACE PROCEDURE reduce_instance_changed(_event eventstore.events2)
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        "name" = _event.payload->>'name'
        , change_date = _event.created_at
        , latest_position = _event.position
        , latest_in_position_order = _event.in_tx_order::INT2
    WHERE 
        id = _event.aggregate_id
        AND (latest_position, latest_in_position_order) < (_event.position, _event.in_tx_order::INT2);
END;
$$;
