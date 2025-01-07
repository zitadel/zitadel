CREATE OR REPLACE FUNCTION reduce_instance_removed("event" eventstore.events2)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instances
    WHERE 
        id = (event).aggregate_id;
END;
$$;
