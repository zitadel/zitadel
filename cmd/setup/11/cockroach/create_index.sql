CREATE INDEX IF NOT EXISTS ca_fill_idx ON eventstore.events (
    event_sequence DESC
    , instance_id
) STORING (
    id
    , creation_date
    , created_at
) WHERE created_at IS NULL; 