CREATE TEMPORARY TABLE IF NOT EXISTS wrong_events (
    instance_id TEXT
    , event_sequence BIGINT
    , current_cd TIMESTAMPTZ
    , next_cd TIMESTAMPTZ
);
