CREATE TABLE projections.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,

    PRIMARY KEY (projection_name)
);

CREATE TABLE projections.current_sequences (
    projection_name TEXT,
    aggregate_type TEXT,
    current_sequence BIGINT,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type)
);

CREATE TABLE projections.failed_events (
    projection_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,
    instance_id TEXT,

    PRIMARY KEY (projection_name, failed_sequence, instance_id)
);
