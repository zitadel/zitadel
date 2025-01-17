CREATE OR REPLACE PROCEDURE reduce_instance_removed(
    _event eventstore.events2
)
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instances
    WHERE 
        id = (_event).aggregate_id;
END;
$$;
