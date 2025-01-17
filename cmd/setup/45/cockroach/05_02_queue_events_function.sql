CREATE OR REPLACE FUNCTION subscriptions.queue_events()
RETURNS TRIGGER 
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO subscriptions.queue (
        subscriber
        , subscriber_name
        , instance_id
        , aggregate_type
        , aggregate_id
        , sequence
        , position
        , in_position_order
        , allow_reduce
        , reduce_function
    )
    SELECT
        s.id
        , s.name
        , (NEW).instance_id
        , (NEW).aggregate_type
        , (NEW).aggregate_id
        , (NEW)."sequence"
        , (NEW).position
        , (NEW).in_tx_order
        , s.allow_reduce AND se.reduce_function IS NOT NULL
        , se.reduce_function
    FROM
        subscriptions.subscribed_events se
    JOIN subscriptions.subscribers s
        ON se.subscriber = s.id
    WHERE
        (se.instance_id IS NULL OR se.instance_id = (NEW).instance_id)
        AND (se.all OR (
            se.aggregate_type = (NEW).aggregate_type
            AND (
                se.event_type IS NULL
                OR se.event_type = (NEW).event_type
            )
        ))
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END;
$$;
