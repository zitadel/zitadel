CREATE OR REPLACE PROCEDURE reduce_instance_domain_added(_event eventstore.events2)
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO instance_domains (
        instance_id
        , domain
        , is_generated
        , creation_date
        , change_date
        , latest_position
        , latest_in_position_order
    ) VALUES (
        _event.aggregate_id
        , _event.payload->>'domain'
        , COALESCE(_event.payload->'generated', 'false')::BOOLEAN
        , _event.created_at
        , _event.created_at
        , _event.position
        , _event.in_tx_order::INT2
    )
    ON CONFLICT (instance_id, domain) DO NOTHING;
END;
$$;
