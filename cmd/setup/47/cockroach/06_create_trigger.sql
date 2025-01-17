CREATE TRIGGER reduce_instance_domain_events
BEFORE INSERT ON subscriptions.queue
FOR EACH ROW
WHEN (NEW).allow_reduce AND (NEW).subscriber_name = 'transactional-instance-domains'
EXECUTE FUNCTION reduce_instance_domain_events();
