CREATE OR REPLACE PROCEDURE reduce_instance_domain_added("event" eventstore.events2)
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
        (event).aggregate_id
        , (event).payload->>'domain'
        , COALESCE((event).payload->'generated', 'false')::BOOLEAN
        , (event).created_at
        , (event).created_at
        , (event).position
        , (event).in_tx_order::INT2
    )
    ON CONFLICT (instance_id, domain) DO NOTHING;
END;
$$;
