CREATE OR REPLACE PROCEDURE reduce_instance_event(
    id UUID
    , instance_id TEXT
    , aggregate_type TEXT
    , aggregate_id TEXT
    , "sequence" INT8
    , event_type TEXT
)
LANGUAGE PLpgSQL
AS $$
BEGIN
    IF event_type = 'instance.added' THEN
        CALL reduce_instance_added(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    ELSIF event_type = 'instance.changed' THEN
        CALL reduce_instance_changed(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    ELSIF event_type = 'instance.removed' THEN
        CALL reduce_instance_removed(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    ELSIF event_type = 'instance.default.language.set' THEN
        CALL reduce_instance_default_language_set(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    ELSIF event_type = 'instance.default.org.set' THEN
        CALL reduce_instance_default_org_set(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    ELSIF event_type = 'instance.iam.project.set' THEN
        CALL reduce_instance_project_set(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    ELSIF event_type = 'instance.iam.console.set' THEN
        CALL reduce_instance_console_set(
            instance_id
            , aggregate_type
            , aggregate_id
            , sequence
        );
    END IF;

    -- DELETE FROM subscriptions.queue WHERE id = id;
END;
$$;
