CREATE SCHEMA IF NOT EXISTS logstore;

GRANT ALL ON ALL TABLES IN SCHEMA logstore TO %[1]s;

CREATE TABLE IF NOT EXISTS logstore.access (
    log_date TIMESTAMPTZ NOT NULL,
    protocol INT NOT NULL,
    request_url TEXT NOT NULL,
    response_status INT NOT NULL,
    request_headers JSONB,
    response_headers JSONB,
    instance_id TEXT NOT NULL,
    project_id TEXT NOT NULL,
    requested_domain TEXT,
    requested_host TEXT
);

CREATE TABLE IF NOT EXISTS logstore.execution (
    log_date TIMESTAMPTZ NOT NULL,
    took INTERVAL,
    message TEXT NOT NULL,
    loglevel INT NOT NULL,
    instance_id TEXT NOT NULL,
    project_id TEXT NOT NULL,
    action_id TEXT NOT NULL,
    metadata JSONB
);

