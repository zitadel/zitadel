CREATE TABLE IF NOT EXISTS subscriptions.subscribers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , name TEXT NOT NULL
    , should_notify BOOLEAN NOT NULL DEFAULT FALSE
    , last_notified_position NUMERIC
    , allow_reduce BOOLEAN NOT NULL DEFAULT FALSE

    , CONSTRAINT name_unique UNIQUE (name)
);
