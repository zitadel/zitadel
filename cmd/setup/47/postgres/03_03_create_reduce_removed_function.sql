CREATE OR REPLACE PROCEDURE reduce_instance_domain_removed("event" eventstore.events2)
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instance_domains
    WHERE 
        instance_id = event.aggregate_id
        AND domain = event.payload->>'domain';
END;
$$;
