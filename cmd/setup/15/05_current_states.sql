CREATE TABLE IF NOT EXISTS projections.current_states (
    projection_name TEXT NOT NULL
    , instance_id TEXT NOT NULL

    , last_updated TIMESTAMPTZ

    , aggregate_id TEXT
    , aggregate_type TEXT
    , "sequence" INT8
    , event_date TIMESTAMPTZ
    , "position" DECIMAL

    , PRIMARY KEY (projection_name, instance_id)
);
CREATE INDEX IF NOT EXISTS cs_instance_id_idx ON projections.current_states (instance_id);