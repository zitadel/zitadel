CREATE OR REPLACE FUNCTION reduce_instance_added(
    id TEXT
    , "name" TEXT 
    , creation_date TIMESTAMPTZ
    , "position" NUMERIC
)
RETURNS VOID
LANGUAGE PLpgSQL
AS $$
BEGIN
    INSERT INTO instances (
        id
        , "name"
        , creation_date
        , change_date
        , latest_position
    ) VALUES (
        id
        , "name"
        , creation_date
        , creation_date
        , "position"
    )
    ON CONFLICT (id) DO NOTHING;
END;
$$;
