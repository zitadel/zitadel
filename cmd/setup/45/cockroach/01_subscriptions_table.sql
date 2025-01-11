CREATE SCHEMA IF NOT EXISTS subscriptions;

CREATE TABLE IF NOT EXISTS subscriptions.subscribers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , name TEXT NOT NULL
    , should_notify BOOLEAN NOT NULL DEFAULT FALSE
    , last_notified_position NUMERIC
    , allow_reduce BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS subscriptions.subscribed_events (
    subscriber UUID NOT NULL
    , instance_id TEXT -- if null susbcription is for all instances
    , "all" BOOLEAN
    , aggregate_type TEXT
    , event_type TEXT

    , reduce_function TEXT -- if null the events are added to the queue

    , CONSTRAINT min_args CHECK (num_nonnulls("all", aggregate_type) = 1)

    , CONSTRAINT subscribers_fk FOREIGN KEY (subscriber) REFERENCES subscriptions.subscribers(id) ON DELETE CASCADE
);