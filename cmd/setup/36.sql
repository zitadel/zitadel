CREATE TABLE IF NOT EXISTS eventstore.snapshots
(
    instance_id text NOT NULL,
    snapshot_type text NOT NULL,
    aggregate_id text NOT NULL,
    "position" numeric NOT NULL,
    change_date timestamptz NOT NULL,
    payload json NOT NULL,

    PRIMARY KEY (instance_id, snapshot_type, aggregate_id)
);
