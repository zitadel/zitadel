CREATE TRIGGER reduce_instance_domain_events
AFTER INSERT ON "queue"
FOR EACH ROW
WHEN (NEW).subscriber = 'transactional-instance-domains'
EXECUTE FUNCTION reduce_instance_domain_events();
