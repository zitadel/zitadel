CREATE OR REPLACE FUNCTION reduce_instance_console_set(
    instance_id TEXT
    , app_id TEXT 
    , client_id TEXT
    , change_date TIMESTAMPTZ
    , "position" NUMERIC
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    UPDATE instances SET
        console_app_id = app_id
        , console_client_id = client_id
        , change_date = change_date
        , latest_position = "position"
    WHERE 
        id = instance_id
        AND latest_position <= "position";
END;
$$;
