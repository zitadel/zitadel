CREATE OR REPLACE PROCEDURE reduce_instance_console_set(
    _instance_id TEXT
    , _aggregate_type TEXT
    , _aggregate_id TEXT
    , _sequence INT8
)
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        console_app_id = event.payload->>'appId'
        , console_client_id = event.payload->>'clientId'
        , change_date = event.created_at
        , latest_position = event.position
        , latest_in_position_order = event.in_tx_order::INT2
    FROM (
        SELECT
            *
        FROM
            eventstore.events2 e
        WHERE
            e.instance_id = _instance_id
            AND e.aggregate_type = _aggregate_type
            AND e.aggregate_id = _aggregate_id
            AND e.sequence = _sequence
    ) AS event
    WHERE 
        id = event.aggregate_id
        AND (latest_position, latest_in_position_order) < (event.position, event.in_tx_order::INT2);
END;
$$;
