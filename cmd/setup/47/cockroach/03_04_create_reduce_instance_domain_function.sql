CREATE OR REPLACE FUNCTION reduce_instance_domain_event(
    _event eventstore.events2
)
RETURNS BOOLEAN
LANGUAGE PLpgSQL
AS $$
BEGIN
    IF (_event).event_type = 'instance.domain.added' THEN
        CALL reduce_instance_domain_added(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.domain.primary.set' THEN
        CALL reduce_instance_domain_primary_set(_event);
        RETURN TRUE;
    ELSIF (_event).event_type = 'instance.domain.removed' THEN
        CALL reduce_instance_domain_removed(_event);
        RETURN TRUE;
    END IF;

    RETURN FALSE;
END;
$$;
