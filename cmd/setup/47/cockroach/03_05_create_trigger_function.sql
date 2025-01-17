CREATE OR REPLACE FUNCTION reduce_instance_domain_events() 
RETURNS TRIGGER
LANGUAGE PLpgSQL
AS $$
DECLARE
    _event eventstore.events2;

    _did_reduce BOOLEAN := FALSE;
BEGIN
    SELECT 
        * 
    INTO
        _event
    FROM 
        eventstore.events2 e
    WHERE 
        e.instance_id = (NEW).instance_id
        AND e.aggregate_type = (NEW).aggregate_type
        AND e.aggregate_id = (NEW).aggregate_id
        AND e.sequence = (NEW).sequence
    ;

    SELECT * INTO _did_reduce FROM reduce_instance_domain_event(_event);

    RAISE NOTICE 'did reduce: %', _did_reduce;

    if _did_reduce THEN
        RETURN NULL;
    END IF;

    RETURN (NEW);
END
$$;
