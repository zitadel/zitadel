CREATE TABLE IF NOT EXISTS logstore.execution (
    log_date TIMESTAMPTZ NOT NULL
    , took INTERVAL
    , message TEXT NOT NULL
    , loglevel INT NOT NULL
    , instance_id TEXT NOT NULL
    , action_id TEXT NOT NULL
    , metadata JSONB
);

CREATE INDEX log_date_desc ON logstore.execution (instance_id, log_date DESC) INCLUDE (took);
