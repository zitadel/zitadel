INSERT INTO eventstore.events (
    instance_id
    , resource_owner
    , aggregate_type
    , aggregate_id
    , aggregate_version

    , editor_user
    , editor_service
    , event_type
    , event_data
    , event_sequence

    , creation_date
    , "position"
    , in_tx_order
) VALUES
    %s
RETURNING creation_date, "position";