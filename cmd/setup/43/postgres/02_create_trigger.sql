DROP TRIGGER IF EXISTS copy_to_outbox ON eventstore.events2;

CREATE OR REPLACE FUNCTION copy_events_to_outbox() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO event_outbox (
        instance_id
        , aggregate_type
        , aggregate_id
        , event_type
        , event_revision
        , created_at
        , payload
        , creator
        , position
        , in_position_order
    ) VALUES (
        (NEW).instance_id
        , (NEW).aggregate_type
        , (NEW).aggregate_id
        , (NEW).event_type
        , (NEW).revision
        , (NEW).created_at
        , (NEW).payload
        , (NEW).creator
        , pg_current_xact_id()::TEXT::NUMERIC
        , (NEW).in_tx_order
    );
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER copy_to_outbox
AFTER INSERT ON eventstore.events2
FOR EACH ROW EXECUTE FUNCTION copy_events_to_outbox();