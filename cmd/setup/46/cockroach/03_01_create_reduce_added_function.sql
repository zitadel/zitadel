CREATE OR REPLACE PROCEDURE reduce_instance_added(
    _instance_id TEXT
    , _aggregate_type TEXT
    , _aggregate_id TEXT
    , _sequence INT8
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
    ) 
    SELECT
        e.aggregate_id
        , e.payload->>'name'
        , e.created_at
        , e.created_at
        , e.position
        , e.in_tx_order::INT2
    FROM
        eventstore.events2 e
    WHERE
        e.instance_id = _instance_id
        AND e.aggregate_type = _aggregate_type
        AND e.aggregate_id = _aggregate_id
        AND e.sequence = _sequence
    ON CONFLICT (id) DO NOTHING;
END;
$$;
