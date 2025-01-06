CREATE TRIGGER copy_to_outbox
AFTER INSERT ON eventstore.events2
FOR EACH ROW EXECUTE FUNCTION copy_events_to_outbox();
