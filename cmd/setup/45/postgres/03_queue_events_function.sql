CREATE OR REPLACE FUNCTION queue_events()
RETURNS TRIGGER 
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO "queue" (
        subscriber
        , instance_id
        , aggregate_type
        , aggregate_id
        , sequence
        , position
        , in_position_order
    ) SELECT 
        subscriber
        , NEW.instance_id
        , NEW.aggregate_type
        , NEW.aggregate_id
        , NEW."sequence"
        , NEW.position
        , NEW.in_tx_order
    FROM
        subscriptions
    WHERE
        (instance_id IS NULL OR NEW.instance_id = instance_id)
        AND ("all" OR (
            aggregate_type = NEW.aggregate_type
            AND (
                event_type IS NULL
                OR NEW.event_type = event_type
            ))
        );
  RETURN NULL;
END;
$$;