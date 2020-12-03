
CREATE SEQUENCE eventstore.event_seq;

GRANT UPDATE ON TABLE eventstore.event_seq TO management;
GRANT UPDATE ON TABLE eventstore.event_seq TO eventstore;
GRANT UPDATE ON TABLE eventstore.event_seq TO adminapi;
GRANT UPDATE ON TABLE eventstore.event_seq TO auth;
GRANT UPDATE ON TABLE eventstore.event_seq TO authz;
GRANT UPDATE ON TABLE eventstore.event_seq TO notification;

SET experimental_enable_hash_sharded_indexes = on;

CREATE TABLE eventstore.events (
    id UUID DEFAULT gen_random_uuid(),
    event_type TEXT,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_version TEXT NOT NULL,
    event_sequence BIGINT NOT NULL DEFAULT nextval('eventstore.event_seq'),
    previous_sequence BIGINT,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    event_data JSONB,
    editor_user TEXT NOT NULL,
    editor_service TEXT NOT NULL,
    resource_owner TEXT NOT NULL,

    CONSTRAINT event_sequence_pk PRIMARY KEY (event_sequence DESC) USING HASH WITH BUCKET_COUNT = 10,
    INDEX agg_type_agg_id (aggregate_type, aggregate_id),
    CONSTRAINT previous_sequence_unique UNIQUE (previous_sequence DESC)
);
ALTER SEQUENCE eventstore.event_seq OWNED BY eventstore.events.event_sequence;
