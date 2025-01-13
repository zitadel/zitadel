CREATE OR REPLACE FUNCTION subscriptions.reduce_events_in_queue(
    subscriber_name TEXT
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
DECLARE
    _stream CURSOR FOR
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
    queued_event RECORD;
BEGIN
    OPEN _stream;
    LOOP
        FETCH _stream INTO queued_event;
        EXIT WHEN NOT FOUND;

        RAISE NOTICE 'Reducing event %', queued_event.event;
        RAISE NOTICE 'execute %', queued_event.reduce_function;

        EXECUTE
            format('CALL %s($1)', queued_event.reduce_function)
        USING 
            queued_event.event;

        DELETE FROM 
            subscriptions.queue q
        WHERE
            q.id = queued_event.id;
            
    END LOOP;
    CLOSE _stream;
END;
$$;