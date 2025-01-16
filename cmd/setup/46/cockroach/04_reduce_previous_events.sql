SELECT subscriptions.reduce_events_in_queue('transactional-instances');

WITH RECURSIVE queued_events(subscriber_id UUID, offset INT) AS (
    SELECT 
        id
        , 0 
    FROM
        subscriptions.subscribers
    WHERE
        name = 'transactional-instances'
    UNION ALL
    WITH queued_event AS (
        SELECT
            q.id AS queue_id
            , e AS "event" 
        FROM
            subscriptions.queue q
        WHERE
            q.subscriber = subscriber_id
        ORDER BY
            position
            , in_tx_order
        LIMIT 1
        OFFSET offset
    ), reduce AS (
        CALL reduce_instance_event((SELECT queue_id FROM queued_event), (SELECT "event" FROM queued_event))
    ) SELECT id, offset+1 FROM queued_events
) select * from queued_events;