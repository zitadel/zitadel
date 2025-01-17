BEGIN;

DECLARE queued_events CURSOR FOR 
    SELECT
        q.id
        , q.instance_id
        , q.aggregate_type
        , q.aggregate_id
        , q.sequence
        , q.event_type
    FROM
        subscriptions.queue q
    WHERE
        q.subscriber = (SELECT id FROM subscriptions.subscribers WHERE name = 'transactional-instances')
    ORDER BY
        q.position
        , q.in_position_order
    ;

CALL reduce_instance_queued_events('queued_events');

COMMIT;
