CREATE TABLE IF NOT EXISTS subscriptions.subscribed_events (
    subscriber UUID NOT NULL
    , instance_id TEXT -- if null susbcription is for all instances
    , "all" BOOLEAN NOT NULL DEFAULT FALSE
    , aggregate_type TEXT
    , event_type TEXT

    , reduce_function TEXT -- if null the events are added to the queue

    -- ,TODO: UNIQUE NULLS NOT DISTINCT (subscriber, instance_id, "all", aggregate_type, event_type, reduce_function)
    , CHECK (CASE WHEN "all" THEN aggregate_type IS NULL ELSE aggregate_type IS NOT NULL END)
    , CONSTRAINT subscribers_fk FOREIGN KEY (subscriber) REFERENCES subscriptions.subscribers(id) ON DELETE CASCADE
);