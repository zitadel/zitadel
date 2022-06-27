CREATE TABLE projections.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (projection_name, instance_id)
);

CREATE TABLE projections.current_sequences (
    projection_name TEXT,
    aggregate_type TEXT,
    current_sequence BIGINT,
    instance_id TEXT NOT NULL,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type, instance_id)
);

CREATE TABLE projections.failed_events (
    projection_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (projection_name, failed_sequence, instance_id)
);
