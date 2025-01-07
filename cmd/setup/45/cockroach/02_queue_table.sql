CREATE TABLE IF NOT EXISTS "queue" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    , subscriber TEXT NOT NULL
    , instance_id TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , sequence INT8 NOT NULL

    , position NUMERIC NOT NULL
    , in_position_order INT2 NOT NULL

    , CONSTRAINT events_fk FOREIGN KEY (instance_id, aggregate_type, aggregate_id, "sequence") REFERENCES eventstore.events2 (instance_id, aggregate_type, aggregate_id, "sequence")
);