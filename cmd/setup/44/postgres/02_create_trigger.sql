DROP TRIGGER IF EXISTS reduce_instance_added ON eventstore.events2;

CREATE OR REPLACE FUNCTION zitadel.reduce_instance_added() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO zitadel.instances (
        id
        , name
        , change_date
        , creation_date
    ) VALUES (
        (NEW).aggregate_id
        , (NEW).payload->>'name'
        , (NEW).created_at
        , (NEW).created_at
    )
    ON CONFLICT (id) DO NOTHING;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reduce_instance_added
AFTER INSERT ON eventstore.events2
FOR EACH ROW
WHEN (NEW).event_type = 'instance.added'
EXECUTE FUNCTION reduce_instance_added();
