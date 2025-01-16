CREATE OR REPLACE PROCEDURE reduce_instance_removed("event" eventstore.events2)
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instances
    WHERE 
        id = (event).aggregate_id;
END;
$$;
