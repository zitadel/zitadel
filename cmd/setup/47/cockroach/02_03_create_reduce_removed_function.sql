CREATE OR REPLACE FUNCTION reduce_instance_domain_removed("event" eventstore.events2)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instance_domains
    WHERE 
        instance_id = (event).aggregate_id
        AND domain = (event).payload->>'domain';
END;
$$;
