CREATE SCHEMA IF NOT EXISTS signals;

CREATE UNLOGGED TABLE IF NOT EXISTS signals.signals (
    instance_id     TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT statement_timestamp(),
    caller_id       TEXT        NOT NULL,
    user_id         TEXT,
    session_id      TEXT,
    fingerprint_id  TEXT,
    stream          TEXT        NOT NULL,
    operation       TEXT        NOT NULL,
    resource        TEXT,
    outcome         TEXT        NOT NULL,
    ip              INET,
    user_agent      TEXT,
    country         TEXT,
    metadata        JSONB
) PARTITION BY RANGE (created_at);

-- Default partition to catch signals before specific partitions are created
CREATE UNLOGGED TABLE IF NOT EXISTS signals.signals_default
    PARTITION OF signals.signals DEFAULT;

-- Risk engine queries by caller within time window
CREATE INDEX IF NOT EXISTS idx_signals_caller
    ON signals.signals (instance_id, caller_id, created_at DESC);

-- Session-scoped queries
CREATE INDEX IF NOT EXISTS idx_signals_session
    ON signals.signals (instance_id, session_id, created_at DESC)
    WHERE session_id IS NOT NULL;

-- Stream-filtered queries
CREATE INDEX IF NOT EXISTS idx_signals_stream
    ON signals.signals (instance_id, caller_id, stream, created_at DESC);

-- Rate limit counters for multi-instance sliding-window rate limiting.
-- UNLOGGED: acceptable to lose on crash (counters reset to zero).
CREATE UNLOGGED TABLE IF NOT EXISTS signals.rate_limit_counters (
    key          TEXT        NOT NULL,
    count        INTEGER     NOT NULL DEFAULT 1,
    window_start TIMESTAMPTZ NOT NULL,
    window_secs  INTEGER     NOT NULL,
    PRIMARY KEY (key)
);

-- Detection rule priority & stop-on-match columns
ALTER TABLE IF EXISTS projections.detection_rules ADD COLUMN IF NOT EXISTS priority BIGINT DEFAULT 0;
ALTER TABLE IF EXISTS projections.detection_rules ADD COLUMN IF NOT EXISTS stop_on_match BOOLEAN DEFAULT FALSE;
