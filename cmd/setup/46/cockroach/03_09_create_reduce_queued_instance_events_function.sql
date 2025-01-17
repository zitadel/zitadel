CREATE OR REPLACE PROCEDURE reduce_instance_queued_events(
    _queued_events REFCURSOR
)
LANGUAGE PLpgSQL
AS $$
DECLARE
    _id UUID;
    _event eventstore.events2;
BEGIN
    LOOP
        FETCH NEXT _queued_events INTO 
            _id
            , _event
        ;
        EXIT WHEN _event IS NULL;

        SELECT reduce_instance_event(_event);

        DELETE FROM subscriptions.queue WHERE id = _id;
    END LOOP;
    CLOSE _queued_events;
END;
$$;
