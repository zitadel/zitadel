CREATE TABLE IF NOT EXISTS instances (
    id TEXT PRIMARY KEY
    , name TEXT NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , latest_position NUMERIC NOT NULL
    , default_org_id TEXT
    , iam_project_id TEXT
    , console_client_id TEXT
    , console_app_id TEXT
    , default_language TEXT
);

-- |     sequence INT8 NOT NULL,