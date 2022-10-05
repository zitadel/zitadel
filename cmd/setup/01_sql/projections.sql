CREATE TABLE projections.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (projection_name, instance_id)
);

CREATE TABLE projections.last_processed_event (
    projection_name TEXT,
    aggregate_type TEXT,
    instance_id TEXT NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT now(),
    event_creation_date TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type, instance_id)
);

CREATE TABLE projections.failed_events (
    projection_name TEXT,
    failed_id UUID,
    failure_count SMALLINT,
    error TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (projection_name, failed_id, instance_id)
);
