CREATE OR REPLACE FUNCTION reduce_instance_default_org_set(
    instance_id TEXT
    , org_id TEXT 
    , change_date TIMESTAMPTZ
    , "position" NUMERIC
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        default_org_id = org_id
        , change_date = change_date
        , latest_position = "position"
    WHERE 
        id = instance_id
        AND latest_position <= "position";
END;
$$;
