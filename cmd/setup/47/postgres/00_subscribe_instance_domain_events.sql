WITH created_subscriber AS (
    INSERT INTO 
        subscriptions.subscribers ("name") 
    VALUES 
        ('transactional-instance-domains')
    ON CONFLICT DO NOTHING
    RETURNING id
)
INSERT INTO subscriptions.subscribed_events (
    subscriber
    , aggregate_type
    , event_type
    , reduce_function
) SELECT
    s.id, (v).*
FROM
    created_subscriber AS s,
    (VALUES
        ('instance', 'instance.domain.added', 'reduce_instance_domain_added')
        , ('instance', 'instance.domain.primary.set', 'reduce_instance_domain_primary_set')
        , ('instance', 'instance.domain.removed', 'reduce_instance_domain_removed')
    )AS v
ON CONFLICT DO NOTHING
;
