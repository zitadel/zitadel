CREATE OR REPLACE FUNCTION reduce_instance_events() 
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

    IF ("event").event_type = 'instance.added' THEN
        SELECT reduce_instance_added("event");
    ELSIF ("event").event_type = 'instance.changed' THEN
        SELECT reduce_instance_changed("event");
    ELSIF ("event").event_type = 'instance.removed' THEN
        SELECT reduce_instance_removed("event");
    ELSIF ("event").event_type = 'instance.default.language.set' THEN
        SELECT reduce_instance_default_language_set("event");
    ELSIF ("event").event_type = 'instance.default.org.set' THEN
        SELECT reduce_instance_default_org_set("event");
    ELSIF ("event").event_type = 'instance.iam.project.set' THEN
        SELECT reduce_instance_project_set("event");
    ELSIF ("event").event_type = 'instance.iam.console.set' THEN
        SELECT reduce_instance_console_set("event");
    END IF;

    DELETE FROM "queue" WHERE id = (NEW).id;
        
    RETURN NULL;
END
$$;
