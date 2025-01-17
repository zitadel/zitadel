CREATE OR REPLACE PROCEDURE reduce_instance_domain_removed(
    _event eventstore.events2
)
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instance_domains
    WHERE 
        instance_id = (_event).aggregate_id
        AND domain = (_event).payload->>'domain';
END;
$$;
