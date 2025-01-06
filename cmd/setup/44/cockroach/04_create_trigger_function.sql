CREATE OR REPLACE FUNCTION reduce_instance_events() 
RETURNS TRIGGER
LANGUAGE PLpgSQL
AS $$
BEGIN
    IF (NEW).event_type = 'instance.added' THEN
        SELECT reduce_instance_added(
            (NEW).aggregate_id
            , (NEW).payload->>'name'::TEXT
            , (NEW).created_at
            , (NEW).position
        );
    ELSIF (NEW).event_type = 'instance.changed' THEN
        SELECT reduce_instance_changed(
            (NEW).aggregate_id
            , (NEW).payload->>'name'::TEXT
            , (NEW).created_at
            , (NEW).position
        );
    ELSIF (NEW).event_type = 'instance.removed' THEN
        SELECT reduce_instance_removed(
            (NEW).aggregate_id
        );
    ELSIF (NEW).event_type = 'instance.default.language.set' THEN
        SELECT reduce_instance_default_language_set(
            (NEW).aggregate_id
            , (NEW).payload->>'language'::TEXT
            , (NEW).created_at
            , (NEW).position
        );
    ELSIF (NEW).event_type = 'instance.default.org.set' THEN
        SELECT reduce_instance_default_org_set(
            (NEW).aggregate_id
            , (NEW).payload->>'orgId'::TEXT
            , (NEW).created_at
            , (NEW).position
        );
    ELSIF (NEW).event_type = 'instance.iam.project.set' THEN
        SELECT reduce_instance_project_set(
            (NEW).aggregate_id
            , (NEW).payload->>'iamProjectId'::TEXT
            , (NEW).created_at
            , (NEW).position
        );
    ELSIF (NEW).event_type = 'instance.iam.console.set' THEN
        SELECT reduce_instance_console_set(
            (NEW).aggregate_id
            , (NEW).payload->>'appId'::TEXT
            , (NEW).payload->>'clientId'::TEXT
            , (NEW).created_at
            , (NEW).position
        );
    END IF;
    RETURN NULL;
END
$$;
