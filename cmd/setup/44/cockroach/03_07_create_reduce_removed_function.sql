CREATE OR REPLACE FUNCTION reduce_instance_removed(
    instance_id TEXT
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instances
    WHERE 
        id = instance_id;
END;
$$;
