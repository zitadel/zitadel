INSERT INTO subscriptions.subscribers ("name") VALUES ('transactional-instances');

INSERT INTO subscriptions.subscribed_events (subscriber, aggregate_type, event_type) 
SELECT
    

VALUES
    ('transactional-instances', 'instance', 'instance.added')
    , ('transactional-instances', 'instance', 'instance.changed')
    , ('transactional-instances', 'instance', 'instance.removed')
    , ('transactional-instances', 'instance', 'instance.default.language.set')
    , ('transactional-instances', 'instance', 'instance.default.org.set')
    , ('transactional-instances', 'instance', 'instance.iam.project.set')
    , ('transactional-instances', 'instance', 'instance.iam.console.set')
;