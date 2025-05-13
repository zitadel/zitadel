CREATE TABLE IF NOT EXISTS aggregates(
    id TEXT NOT NULL
    , type TEXT NOT NULL
    , instance_id TEXT NOT NULL
    
    , current_sequence INT NOT NULL DEFAULT 0

    , PRIMARY KEY (instance_id, type, id)
);

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()

    -- object type that the event is related to
    , aggregate TEXT NOT NULL
    -- id of the object that the event is related to
    , aggregate_id TEXT NOT NULL
    , instance_id TEXT NOT NULL
    
    -- time the event was created
    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    -- user that created the event
    , creator TEXT
    -- type of the event
    , type TEXT NOT NULL
    -- version of the event
    , revision SMALLINT NOT NULL
    -- changed fields or NULL
    , payload JSONB

    , position NUMERIC NOT NULL DEFAULT pg_current_xact_id()::TEXT::NUMERIC
    , in_position_order INT2 NOT NULL

    , FOREIGN KEY (instance_id, aggregate, aggregate_id) REFERENCES aggregates(instance_id, type, id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS instances(
    id TEXT
    , name TEXT NOT NULL
    , created_at TIMESTAMPTZ NOT NULL
    , updated_at TIMESTAMPTZ NOT NULL

    , default_org_id TEXT
    , iam_project_id TEXT
    , console_client_id TEXT
    , console_app_id TEXT
    , default_language VARCHAR(10)

    , PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS instance_domains(
    instance_id TEXT NOT NULL
    , domain TEXT NOT NULL
    , is_primary BOOLEAN NOT NULL DEFAULT FALSE
    , is_verified BOOLEAN NOT NULL DEFAULT FALSE
    
    , PRIMARY KEY (instance_id, domain)
    , FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS instance_domain_search_idx ON instance_domains (domain);