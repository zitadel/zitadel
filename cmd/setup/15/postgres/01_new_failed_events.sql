CREATE TABLE IF NOT EXISTS projections.failed_events2 (
    projection_name TEXT NOT NULL
    , instance_id TEXT NOT NULL

    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , event_creation_date TIMESTAMPTZ NOT NULL
    , failed_sequence INT8 NOT NULL

    , failure_count INT2 NULL DEFAULT 0
    , error TEXT
    , last_failed TIMESTAMPTZ
    
    , PRIMARY KEY (projection_name, instance_id, aggregate_type, aggregate_id, failed_sequence)
);
CREATE INDEX IF NOT EXISTS fe2_instance_id_idx on projections.failed_events2 (instance_id);