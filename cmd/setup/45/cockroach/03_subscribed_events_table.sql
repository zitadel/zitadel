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