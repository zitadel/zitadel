CREATE OR REPLACE FUNCTION reduce_instance_events() 
RETURNS TRIGGER
LANGUAGE PLpgSQL
AS $$
DECLARE
    id UUID := (NEW).id;
    instance_id TEXT := (NEW).instance_id;
    aggregate_type TEXT := (NEW).aggregate_type;
    aggregate_id TEXT := (NEW).aggregate_id;
    "sequence" INT8 := (NEW)."sequence";
    event_type TEXT := (NEW).event_type;
BEGIN
    CALL reduce_instance_event(
        id
        , instance_id
        , aggregate_type
        , aggregate_id
        , "sequence"
        , event_type
    );
    RETURN (NEW);
END
$$;
