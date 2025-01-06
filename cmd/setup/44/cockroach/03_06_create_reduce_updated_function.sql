CREATE OR REPLACE FUNCTION reduce_instance_changed(
    instance_id TEXT
    , "name" TEXT 
    , change_date TIMESTAMPTZ
    , "position" NUMERIC
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        "name" = $2
        , change_date = change_date
        , latest_position = "position"
    WHERE 
        id = instance_id
        AND latest_position <= "position";
END;
$$;
