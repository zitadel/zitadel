BEGIN;

UPDATE subscriptions.subscribers
SET allow_reduce = TRUE
WHERE name = 'transactional-instances';

SELECT subscriptions.reduce_events_in_queue('transactional-instances');

COMMIT;
