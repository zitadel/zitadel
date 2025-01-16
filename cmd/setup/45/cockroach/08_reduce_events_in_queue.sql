CREATE OR REPLACE FUNCTION subscriptions.reduce_events_in_queue(
    subscriber_name TEXT

    _queued_events REFCURSOR
)
RETURNS REFCURSOR
LANGUAGE PLpgSQL
AS $$
DECLARE
    _stream CURSOR FOR
        SELECT
            q.id
            , e AS "event"
        FROM
            subscriptions.subscribers s
        JOIN subscriptions.queue q
            ON q.subscriber = s.id
        JOIN 
            eventstore.events2 e
            ON e.instance_id = q.instance_id
            AND e.aggregate_type = q.aggregate_type
            AND e.aggregate_id = q.aggregate_id
            AND e."sequence" = q.sequence
        WHERE
            s.name = subscriber_name
        ORDER BY 
            q.position
            , q.in_position_order;
    queued_event RECORD;
BEGIN
    OPEN _queued_events FOR
        SELECT
            q.id
            , q.reduce_function AS reduce_function
            , e AS "event"
        FROM
            subscriptions.subscribers s
        JOIN subscriptions.queue q
            ON q.subscriber = s.id
        JOIN 
            eventstore.events2 e
            ON e.instance_id = q.instance_id
            AND e.aggregate_type = q.aggregate_type
            AND e.aggregate_id = q.aggregate_id
            AND e."sequence" = q.sequence
        WHERE
            s.name = subscriber_name
        ORDER BY 
            q.position
            , q.in_position_order;
END;
$$;