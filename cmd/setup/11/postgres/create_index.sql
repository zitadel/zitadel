CREATE INDEX IF NOT EXISTS ca_fill_idx ON eventstore.events (
    event_sequence DESC
    , instance_id
) WHERE created_at IS NULL; 