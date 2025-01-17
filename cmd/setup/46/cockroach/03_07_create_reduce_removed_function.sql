CREATE OR REPLACE PROCEDURE reduce_instance_removed(
    _instance_id TEXT
    , _aggregate_type TEXT
    , _aggregate_id TEXT
    , _sequence INT8
)
LANGUAGE PLpgSQL
AS $$
BEGIN
    DELETE FROM 
        instances
    WHERE 
        id = _aggregate_id;
END;
$$;
