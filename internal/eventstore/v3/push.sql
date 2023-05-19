INSERT INTO eventstore.events (
    instance_id
    , resource_owner
    , aggregate_type
    , aggregate_id
    , aggregate_version

    , editor_user
    , event_type
    , event_data
    , event_sequence
) VALUES
    %s
RETURNING creation_date;