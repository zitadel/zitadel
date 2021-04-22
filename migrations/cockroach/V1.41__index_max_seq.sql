use eventstore;
CREATE INDEX IF NOT EXISTS max_sequence on eventstore.events (aggregate_type, aggregate_id, event_sequence DESC);
