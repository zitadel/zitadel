WITH active_instances AS (
    SELECT
        instance_id
    FROM
        eventstore.events2
    WHERE
        aggregate_type = 'instance'
        AND event_type = 'instance.added'
        AND instance_id NOT IN (
            SELECT
                instance_id
            FROM
                eventstore.events2
            WHERE
                aggregate_type = 'instance'
                AND event_type = 'instance.removed'
        )
)
INSERT INTO "queue" (
    subscriber
    , instance_id
    , aggregate_type
    , aggregate_id
    , sequence
    , position
    , in_position_order
) SELECT
    'transactional-instance-domains'
    , e.instance_id
    , e.aggregate_type
    , e.aggregate_id
    , e."sequence"
    , e.position
    , e.in_tx_order
FROM
    eventstore.events2 e
WHERE
    e.instance_id IN (SELECT instance_id FROM active_instances)
    AND e.aggregate_type = 'instance'
    AND e.event_type IN ('instance.domain.added', 'instance.domain.primary.set', 'instance.domain.removed')
    AND e.aggregate_id IN (SELECT instance_id FROM active_instances);
