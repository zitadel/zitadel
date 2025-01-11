CREATE TABLE IF NOT EXISTS subscriptions."queue" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , subscriber UUID NOT NULL
    , instance_id TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , sequence INT8 NOT NULL

    , position NUMERIC NOT NULL
    , in_position_order INT2 NOT NULL

    , allow_reduce BOOLEAN NOT NULL DEFAULT FALSE
    , reduce_function TEXT

    , CONSTRAINT subscribers_fk FOREIGN KEY (subscriber) REFERENCES subscriptions.subscribers(id) ON DELETE CASCADE
    , CONSTRAINT events_fk FOREIGN KEY (instance_id, aggregate_type, aggregate_id, "sequence") REFERENCES eventstore.events2 (instance_id, aggregate_type, aggregate_id, "sequence")

    , UNIQUE (subscriber, instance_id, aggregate_type, aggregate_id, "sequence")
);
