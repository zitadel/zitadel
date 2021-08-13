CREATE DATABASE zitadel;
GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE zitadel TO queries;
use zitadel;

CREATE SCHEMA zitadel.projections AUTHORIZATION queries;

CREATE TABLE zitadel.projections.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,

    PRIMARY KEY (projection_name)
);

CREATE TABLE zitadel.projections.current_sequences (
    projection_name TEXT,
    aggregate_type TEXT,
    current_sequence BIGINT,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type)
);

CREATE TABLE zitadel.projections.failed_events (
    projection_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,

    PRIMARY KEY (projection_name, failed_sequence)
);
