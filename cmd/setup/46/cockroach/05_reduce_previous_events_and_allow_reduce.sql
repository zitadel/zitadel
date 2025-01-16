BEGIN;

UPDATE subscriptions.subscribers
SET allow_reduce = TRUE
WHERE name = 'transactional-instances';

SELECT subscriptions.reduce_events_in_queue('transactional-instances');

CREATE TRIGGER reduce_instance_events
AFTER INSERT ON "queue"
FOR EACH ROW
WHEN (NEW).subscriber = 'transactional-instances'
EXECUTE FUNCTION reduce_instance_events();

COMMIT;
