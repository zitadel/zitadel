CREATE DATABASE eventstore;

CREATE USER eventstore;
GRANT SELECT, INSERT, UPDATE ON DATABASE eventstore TO eventstore;

CREATE SEQUENCE eventstore.event_seq;

CREATE TABLE eventstore.events (
    event_type TEXT,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_version TEXT NOT NULL,
    event_sequence BIGINT NOT NULL DEFAULT nextval('eventstore.event_seq'),
    previous_sequence BIGINT UNIQUE,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    event_data JSONB,
    editor_user TEXT NOT NULL, 
    editor_service TEXT NOT NULL,
    resource_owner TEXT NOT NULL,

    PRIMARY KEY (event_sequence DESC),
    INDEX (aggregate_type, aggregate_id)
);
