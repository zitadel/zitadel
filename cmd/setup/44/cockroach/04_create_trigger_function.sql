CREATE OR REPLACE FUNCTION reduce_instance_events() 
RETURNS TRIGGER
LANGUAGE PLpgSQL
AS $$
BEGIN
    IF (NEW).event_type = 'instance.added' THEN
        SELECT reduce_instance_added(NEW::eventstore.events2);
    -- ELSIF (NEW).event_type = 'instance.changed' THEN
    --     SELECT reduce_instance_changed(NEW::eventstore.events2);
    -- ELSIF (NEW).event_type = 'instance.removed' THEN
    --     SELECT reduce_instance_removed(NEW::eventstore.events2);
    -- ELSIF (NEW).event_type = 'instance.default.language.set' THEN
    --     SELECT reduce_instance_default_language_set(NEW::eventstore.events2);
    -- ELSIF (NEW).event_type = 'instance.default.org.set' THEN
    --     SELECT reduce_instance_default_org_set(NEW::eventstore.events2);
    -- ELSIF (NEW).event_type = 'instance.iam.project.set' THEN
    --     SELECT reduce_instance_project_set(NEW::eventstore.events2);
    -- ELSIF (NEW).event_type = 'instance.iam.console.set' THEN
    --     SELECT reduce_instance_console_set(NEW::eventstore.events2);
    END IF;
    RETURN NULL;
END
$$;
