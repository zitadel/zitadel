CREATE TABLE eventstore.events (
    event_type TEXT,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_version TEXT NOT NULL,
    event_sequence INTEGER,
    previous_sequence BIGINT,
    creation_date TIMESTAMPT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event_data JSONB,
    editor_user TEXT NOT NULL, 
    editor_service TEXT NOT NULL,
    resource_owner TEXT NOT NULL,

    CONSTRAINT event_sequence_pk PRIMARY KEY (event_sequence DESC),
    CONSTRAINT previous_sequence_unique UNIQUE (previous_sequence DESC)
);

CREATE INDEX eventstore.agg_type_agg_id ON events (aggregate_type, aggregate_id);
