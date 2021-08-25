BEGIN;

SET experimental_enable_hash_sharded_indexes = true;

CREATE TABLE eventstore.events_new (
    id UUID DEFAULT gen_random_uuid(),
    event_type TEXT,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_version TEXT NOT NULL,
    event_sequence BIGINT NOT NULL DEFAULT nextval('eventstore.event_seq'),
    previous_aggregate_sequence BIGINT,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    event_data JSONB,
    editor_user TEXT NOT NULL, 
    editor_service TEXT NOT NULL,
    resource_owner TEXT NOT NULL,
    previous_aggregate_type_sequence BIGINT,

    CONSTRAINT event_sequence_pk PRIMARY KEY (event_sequence DESC) USING HASH WITH BUCKET_COUNT = 10,
    CONSTRAINT previous_sequence_unique UNIQUE (previous_aggregate_sequence DESC),
    INDEX agg_type_agg_id (aggregate_type, aggregate_id),
    INDEX max_sequence (aggregate_type, aggregate_id, event_sequence DESC),
    INDEX default_event_query (aggregate_type, aggregate_id, event_type, resource_owner),
    INDEX agg_type (aggregate_type)
);

COMMIT;
BEGIN;

INSERT INTO eventstore.events_new(
    id,
    event_type,
    aggregate_type,
    aggregate_id,
    aggregate_version,
    event_sequence,
    previous_aggregate_sequence,
    creation_date,
    event_data,
    editor_user,
    editor_service,
    resource_owner,
    previous_aggregate_type_sequence
) SELECT 
    id,
    event_type,
    aggregate_type,
    aggregate_id,
    aggregate_version,
    event_sequence,
    previous_aggregate_sequence,
    creation_date,
    event_data,
    editor_user,
    editor_service,
    resource_owner,
    LAG(event_sequence) 
        OVER (
            PARTITION BY aggregate_type 
            ORDER BY event_sequence
        ) as previous_aggregate_type_sequence
    FROM eventstore.events
    ORDER BY event_sequence;

COMMIT;
BEGIN;

ALTER TABLE eventstore.events RENAME TO events_old;
ALTER SEQUENCE eventstore.event_seq OWNED BY eventstore.events_new.previous_aggregate_sequence;

COMMIT;
BEGIN;

ALTER TABLE eventstore.events_new RENAME TO events;

COMMIT;
BEGIN;

INSERT INTO eventstore.events (
    id,
    event_type,
    aggregate_type,
    aggregate_id,
    aggregate_version,
    event_sequence,
    previous_aggregate_sequence,
    creation_date,
    event_data,
    editor_user,
    editor_service,
    resource_owner,
    previous_aggregate_type_sequence
) (SELECT 
    id,
    event_type,
    aggregate_type,
    aggregate_id,
    aggregate_version,
    event_sequence,
    previous_aggregate_sequence,
    creation_date,
    event_data,
    editor_user,
    editor_service,
    resource_owner,
    previous_aggregate_type_sequence
FROM eventstore.events_old
WHERE event_sequence > (SELECT MAX(event_sequence) FROM eventstore.events));

COMMIT;
