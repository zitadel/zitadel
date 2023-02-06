INSERT INTO zitadel.eventstore_v3.events (
    aggregate_id
    , aggregate_type
    , owner
    , instance_id
    , user_id
    , service
    , event_type
    , event_version
    , payload
    , creation_date
) VALUES %s RETURNING creation_date