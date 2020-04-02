BEGIN;

CREATE DATABASE eventstore;

COMMIT;


BEGIN;

CREATE USER eventstore;

GRANT SELECT, INSERT, UPDATE ON DATABASE eventstore TO eventstore;

COMMIT;

BEGIN;

CREATE TABLE eventstore.events (
    id UUID DEFAULT gen_random_uuid(),
    
    event_type TEXT,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    aggregate_version TEXT NOT NULL,
    event_sequence BIGSERIAL,
    previous_sequence BIGINT UNIQUE,
    creation_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    event_data JSONB,
    modifier_user TEXT NOT NULL, 
    modifier_service TEXT NOT NULL,
    modifier_tenant TEXT NOT NULL,
    resource_owner TEXT NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE eventstore.locks (
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    until TIMESTAMPTZ,
    UNIQUE (aggregate_type, aggregate_id)
);

COMMIT;
