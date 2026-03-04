INSERT INTO eventstore.events2 (
    instance_id
    , "owner"
    , aggregate_type
    , aggregate_id
    , revision

    , creator
    , event_type
    , payload
    , "sequence"
    , created_at

    , "position"
    , in_tx_order
    , written_by_v3
) VALUES
    %s
RETURNING created_at, "position";