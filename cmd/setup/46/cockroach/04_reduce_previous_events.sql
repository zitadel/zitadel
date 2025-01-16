BEGIN;

DECLARE queued_events CURSOR FOR 
    SELECT
        q.id
        , e AS "event"
    FROM
        subscriptions.queue q
    JOIN
        eventstore.events2 e
        ON e.instance_id = q.instance_id
        AND e.aggregate_type = q.aggregate_type
        AND e.aggregate_id = q.aggregate_id
        AND e."sequence" = q.sequence
    WHERE
        q.subscriber = (SELECT id FROM subscriptions.subscribers WHERE name = 'transactional-instances');

OPEN queued_events;


CLOSE queued_events;