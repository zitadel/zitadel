CREATE SCHEMA IF NOT EXISTS logstore;

GRANT ALL ON ALL TABLES IN SCHEMA logstore TO %[1]s;

CREATE TABLE IF NOT EXISTS logstore.access (
    ts TIMESTAMPTZ,
    protocol INT,
    request_url TEXT,
    response_status INT,
    request_headers JSONB,
    response_headers JSONB,
    instance_id      TEXT,
    project_id       TEXT,
    requested_domain TEXT,
    requested_host   TEXT
);

CREATE TABLE IF NOT EXISTS logstore.execution (
    ts TIMESTAMPTZ,
    started TIMESTAMPTZ,
    message TEXT,
    loglevel INT,
    instance_id      TEXT,
    project_id TEXT,
    action_id TEXT,
    metadata JSONB
);

