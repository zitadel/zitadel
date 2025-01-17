WITH created_subscriber AS (
    INSERT INTO 
        subscriptions.subscribers ("name") 
    VALUES 
        ('transactional-instances')
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
        ('instance', 'instance.added', 'reduce_instance_added')
        , ('instance', 'instance.changed', 'reduce_instance_changed')
        , ('instance', 'instance.removed', 'reduce_instance_removed')
        , ('instance', 'instance.default.language.set', 'reduce_instance_default_language_set')
        , ('instance', 'instance.default.org.set', 'reduce_instance_default_org_set')
        , ('instance', 'instance.iam.project.set', 'reduce_instance_project_set')
        , ('instance', 'instance.iam.console.set', 'reduce_instance_console_set')
    )AS v
;