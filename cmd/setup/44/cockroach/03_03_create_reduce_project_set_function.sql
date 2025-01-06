CREATE OR REPLACE FUNCTION reduce_instance_project_set(
    instance_id TEXT
    , project_id TEXT 
    , change_date TIMESTAMPTZ
    , "position" NUMERIC
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        iam_project_id = project_id
        , change_date = change_date
        , latest_position = "position"
    WHERE 
        id = instance_id
        AND latest_position <= "position";
END;
$$;
