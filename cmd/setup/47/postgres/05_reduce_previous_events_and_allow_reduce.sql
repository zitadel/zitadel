BEGIN;

UPDATE subscriptions.subscribers
SET allow_reduce = TRUE
WHERE name = 'transactional-instance-domains';

SELECT subscriptions.reduce_events_in_queue('transactional-instance-domains');

COMMIT;
