INSERT INTO zitadel.eventstore_v5.events (
    aggregate_id
    , aggregate_type
    , owner
    , instance_id
    , user_id
    , service
    , event_type
    , event_version
    , payload
    , sequence
) VALUES %s