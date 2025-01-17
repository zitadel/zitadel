CREATE TRIGGER reduce_instance_events
BEFORE INSERT ON subscriptions.queue
FOR EACH ROW
WHEN (NEW).allow_reduce AND (NEW).subscriber_name = 'transactional-instances'
EXECUTE FUNCTION reduce_instance_events();
