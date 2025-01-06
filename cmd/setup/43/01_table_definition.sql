CREATE TABLE IF NOT EXISTS event_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , instance_id TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , event_type TEXT NOT NULL
    , event_revision INT2 NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT TRANSACTION_TIMESTAMP()
    , payload JSONB NULL
    , creator TEXT NOT NULL
    , position NUMERIC NOT NULL
    , in_position_order INT4 NOT NULL
);