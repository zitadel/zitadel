CREATE OR REPLACE PROCEDURE reduce_instance_default_language_set("event" eventstore.events2)
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        default_language = event.payload->>'language'
        , change_date = event.created_at
        , latest_position = event.position
        , latest_in_position_order = event.in_tx_order::INT2
    WHERE 
        id = event.aggregate_id
        AND (latest_position, latest_in_position_order) < (event.position, event.in_tx_order::INT2);
END;
$$;
