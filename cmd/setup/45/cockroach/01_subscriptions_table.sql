CREATE TABLE IF NOT EXISTS subscriptions (
    subscriber TEXT NOT NULL
    , instance_id TEXT -- if null susbcription is for all instances
    , "all" BOOLEAN
    , aggregate_type TEXT
    , event_type TEXT

    , CONSTRAINT min_args CHECK ("all" OR aggregate_type IS NOT NULL)
);