-- sql statements executed to init zitadel:

-- user
CREATE USER zitadel WITH PASSWORD 'zitadel';

-- database
CREATE DATABASE zitadel;

-- grant user
GRANT ALL ON DATABASE zitadel TO zitadel;

\c zitadel;
SET ROLE zitadel;

--schemas
CREATE SCHEMA eventstore;
CREATE SCHEMA projections;
CREATE SCHEMA system;

-- encryption keys
CREATE TABLE system.encryption_keys (
	id TEXT NOT NULL
	, key TEXT NOT NULL

	, PRIMARY KEY (id)
);

-- events table
CREATE TABLE eventstore.events (
	id UUID DEFAULT gen_random_uuid()
	, event_type TEXT NOT NULL
	, aggregate_type TEXT NOT NULL
	, aggregate_id TEXT NOT NULL
	, aggregate_version TEXT NOT NULL
	, event_sequence BIGINT NOT NULL
	, previous_aggregate_sequence BIGINT
	, previous_aggregate_type_sequence INT8
	, creation_date TIMESTAMPTZ NOT NULL DEFAULT now()
	, event_data JSONB
	, editor_user TEXT NOT NULL 
	, editor_service TEXT NOT NULL
	, resource_owner TEXT NOT NULL
	, instance_id TEXT NOT NULL

	, PRIMARY KEY (event_sequence, instance_id)
	, CONSTRAINT previous_sequence_unique UNIQUE(previous_aggregate_sequence, instance_id)
	, CONSTRAINT prev_agg_type_seq_unique UNIQUE(previous_aggregate_type_sequence, instance_id)
);

CREATE INDEX agg_type_agg_id ON eventstore.events (aggregate_type, aggregate_id, instance_id);
CREATE INDEX agg_type ON eventstore.events (aggregate_type, instance_id);
CREATE INDEX agg_type_seq ON eventstore.events (aggregate_type, event_sequence DESC, instance_id);
CREATE INDEX max_sequence ON eventstore.events (aggregate_type, aggregate_id, event_sequence DESC, instance_id);

-- unique constraints
CREATE TABLE eventstore.unique_constraints (
    instance_id TEXT,
    unique_type TEXT,
    unique_field TEXT,
    PRIMARY KEY (instance_id, unique_type, unique_field)
);

-- system sequence
CREATE SEQUENCE eventstore.system_seq;
