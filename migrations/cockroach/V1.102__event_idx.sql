CREATE INDEX agg_type_seq ON eventstore.events 
    (aggregate_type, event_sequence DESC) 
STORING (
    id
    , event_type
    , aggregate_id
    , aggregate_version
    , previous_aggregate_sequence
    , creation_date
    , event_data
    , editor_user
    , editor_service
    , resource_owner
    , previous_aggregate_type_sequence
);