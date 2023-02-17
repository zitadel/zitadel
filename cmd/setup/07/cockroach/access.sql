CREATE TABLE IF NOT EXISTS logstore.access (
    log_date TIMESTAMPTZ NOT NULL
    , protocol INT NOT NULL
    , request_url TEXT NOT NULL
    , response_status INT NOT NULL
    , request_headers JSONB
    , response_headers JSONB
    , instance_id TEXT NOT NULL
    , project_id TEXT NOT NULL
    , requested_domain TEXT
    , requested_host TEXT

    , INDEX protocol_date_desc (instance_id, protocol, log_date DESC) STORING (request_url, response_status, request_headers)
);
