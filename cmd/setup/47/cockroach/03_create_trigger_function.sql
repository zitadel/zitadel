CREATE OR REPLACE FUNCTION reduce_instance_domain_events() 
RETURNS TRIGGER
LANGUAGE PLpgSQL
AS $$
DECLARE
    "event" eventstore.events2;
BEGIN
    SELECT 
        * 
    INTO
        event
    FROM 
        eventstore.events2 e
    WHERE 
        e.instance_id = (NEW).instance_id
        AND e.aggregate_type = (NEW).aggregate_type
        AND e.aggregate_id = (NEW).aggregate_id
        AND e."sequence" = (NEW)."sequence"
    ;

    IF ("event").event_type = 'instance.domain.added' THEN
        SELECT reduce_instance_domain_added("event");
    ELSIF ("event").event_type = 'instance.domain.primary.set' THEN
        SELECT reduce_instance_domain_primary_set("event");
    ELSIF ("event").event_type = 'instance.domain.removed' THEN
        SELECT reduce_instance_domain_removed("event");
    END IF;

    DELETE FROM "queue" WHERE id = (NEW).id;
        
    RETURN NULL;
END
$$;
