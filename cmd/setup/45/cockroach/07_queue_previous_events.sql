CREATE OR REPLACE FUNCTION subscriptions.queue_previous_events(
    subscriber_name TEXT
    , max_position NUMERIC
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO subscriptions.queue (
        subscriber
        , instance_id
        , aggregate_type
        , aggregate_id
        , sequence
        , position
        , in_position_order
        , reduce_function
    ) 
    SELECT
        s.id
        , e.instance_id
        , e.aggregate_type
        , e.aggregate_id
        , e."sequence"
        , e.position
        , e.in_tx_order
        , se.reduce_function
    FROM
        subscriptions.subscribers s
    JOIN subscriptions.subscribed_events se
        ON se.subscriber = s.id
    JOIN eventstore.events2 e
        ON (se.instance_id IS NULL OR se.instance_id = e.instance_id)
        AND (se.all OR (
            se.aggregate_type = e.aggregate_type
            AND (
                se.event_type IS NULL
                OR se.event_type = e.event_type
            )))
        AND (max_position IS NULL OR e.position < max_position)
    WHERE
        s.name = subscriber_name
    ON CONFLICT DO NOTHING;
END;
$$;
