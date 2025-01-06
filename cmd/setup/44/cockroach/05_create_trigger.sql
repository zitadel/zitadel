CREATE TRIGGER reduce_instance_events
AFTER INSERT ON eventstore.events2
FOR EACH ROW
WHEN (NEW).aggregate_type = 'instance'
EXECUTE FUNCTION reduce_instance_events();
