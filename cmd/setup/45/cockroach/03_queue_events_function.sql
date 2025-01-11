CREATE OR REPLACE FUNCTION subscriptions.queue_events()
RETURNS TRIGGER 
LANGUAGE PLpgSQL
AS $$
DECLARE
    _cursor CURSOR FOR 
        SELECT
            s.id
            , s.allow_reduce AND se.reduce_function IS NOT NULL AS use_reduce_function
            , se.reduce_function
        FROM
            subscriptions.subscribers s
        JOIN subscriptions.subscribed_events se
          ON se.subscriber = s.id
         AND (se.instance_id IS NULL OR se.instance_id = NEW.instance_id)
         AND (se."all" OR (
            se.aggregate_type = NEW.aggregate_type
            AND (
                se.event_type IS NULL
                OR NEW.event_type = se.event_type
            ))
        );
    subscriber_id UUID;
    use_reduce_function BOOLEAN;
    reduce_function TEXT;
    "event" eventstore.events2;
BEGIN
    LOOP
        FETCH _cursor INTO subscriber_id, use_reduce_function, reduce_function;
        EXIT WHEN NOT FOUND;

        IF use_reduce_function THEN
            SELECT
                *
            INTO
                "event"
            FROM
                eventstore.events2 e
            WHERE
                e.instance_id = NEW.instance_id
                AND e.aggregate_type = NEW.aggregate_type
                AND e.aggregate_id = NEW.aggregate_id
                AND e."sequence" = NEW."sequence";
            EXECUTE 
                format('%s($1)', subscriber.reduce_function)
            USING
                "event";
        ELSE
            INSERT INTO subscriptions.queue VALUES(
                subscriber.id
                , NEW.instance_id
                , NEW.aggregate_type
                , NEW.aggregate_id
                , NEW."sequence"
                , NEW.position
                , NEW.in_tx_order
            );
        END IF;
    END LOOP;
    RETURN NEW;
END;
$$;