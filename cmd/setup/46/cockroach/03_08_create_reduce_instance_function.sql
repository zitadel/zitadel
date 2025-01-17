CREATE OR REPLACE FUNCTION reduce_instance_event(
    _event eventstore.events2
)
RETURNS BOOLEAN
LANGUAGE PLpgSQL
AS $$
BEGIN
    IF (_event).event_type = 'instance.added' THEN
        CALL reduce_instance_added(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.changed' THEN
        CALL reduce_instance_changed(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.removed' THEN
        CALL reduce_instance_removed(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.default.language.set' THEN
        CALL reduce_instance_default_language_set(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.default.org.set' THEN
        CALL reduce_instance_default_org_set(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.iam.project.set' THEN
        CALL reduce_instance_project_set(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.iam.console.set' THEN
        CALL reduce_instance_console_set(_event);
        RETURN TRUE;
    END IF;
    RETURN FALSE;
END;
$$;
