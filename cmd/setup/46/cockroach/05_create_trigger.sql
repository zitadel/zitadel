CREATE TRIGGER reduce_instance_events
AFTER INSERT ON "queue"
FOR EACH ROW
WHEN (NEW).subscriber = 'transactional-instances'
EXECUTE FUNCTION reduce_instance_events();
