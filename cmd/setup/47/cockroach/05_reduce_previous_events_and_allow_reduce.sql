BEGIN;

UPDATE subscriptions.subscribers
SET allow_reduce = TRUE
WHERE name = 'transactional-instance-domains';

DECLARE queued_events CURSOR FOR 
    SELECT
        q.id
        , e
    FROM
        subscriptions.queue q
    JOIN
        eventstore.events2 e
    ON
        q.instance_id = e.instance_id
        AND q.aggregate_type = e.aggregate_type
        AND q.aggregate_id = e.aggregate_id
        AND q.sequence = e.sequence
    WHERE
        q.subscriber = (SELECT id FROM subscriptions.subscribers WHERE name = 'transactional-instance-domains')
    ORDER BY
        q.position
        , q.in_position_order
    ;

CALL reduce_instance_queued_events('queued_events');

COMMIT;
