CREATE OR REPLACE PROCEDURE reduce_instance_queued_events(
    _queued_events REFCURSOR
)
LANGUAGE PLpgSQL
AS $$
DECLARE
    id UUID;
    instance_id TEXT;
    aggregate_type TEXT;
    aggregate_id TEXT;
    event_type TEXT;
    "sequence" INT8;
BEGIN
    LOOP
        FETCH NEXT _queued_events INTO 
            id
            , instance_id
            , aggregate_type
            , aggregate_id
            , event_type
            , sequence;
        
        EXIT WHEN id IS NULL;

        CALL reduce_instance_event(
            id
            , instance_id
            , aggregate_type
            , aggregate_id
            , sequence
            , event_type
        );
    END LOOP;
    CLOSE _queued_events;
END;
$$;
