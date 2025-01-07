CREATE OR REPLACE FUNCTION reduce_instance_default_org_set("event" eventstore.events2)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        default_org_id = (event).payload->>'orgId'
        , change_date = (event).created_at
        , latest_position = (event).position
        , latest_in_position_order = (event).in_tx_order::INT2
    WHERE 
        id = (event).aggregate_id
        AND (latest_position, latest_in_position_order) < ((event).position, (event).in_tx_order::INT2);
END;
$$;
