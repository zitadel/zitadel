CREATE DATABASE IF NOT EXISTS zitadel;

CREATE SCHEMA IF NOT EXISTS zitadel.eventstore_v4;

CREATE TABLE IF NOT EXISTS zitadel.eventstore_v4.sequences (
    aggregate_id STRING NOT NULL
    , instance_id STRING NOT NULL
    , sequence INT8 NOT NULL
    , region STRING

    , PRIMARY KEY (instance_id, aggregate_id)
);

CREATE TABLE IF NOT EXISTS zitadel.eventstore_v4.events (
    -- aggregate information
    aggregate_id STRING NOT NULL
    , aggregate_type STRING NOT NULL
    , instance_id STRING NOT NULL
    , owner STRING NOT NULL
    -- editor metadata
    , user_id STRING NOT NULL
    , service STRING NOT NULL
    -- event metadata
    , event_type STRING NOT NULL
    , event_version SMALLINT NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL DEFAULT now()
    , payload JSONB
    , region STRING
    , sequence INT8 NOT NULL 

    -- , INDEX agg_fk (aggregate_id, instance_id)
    -- , CONSTRAINT fk_agg FOREIGN KEY (instance_id, aggregate_id) REFERENCES zitadel.eventstore_v4.aggregates (instance_id, id)
    
    , CONSTRAINT agg_seq UNIQUE(instance_id, aggregate_id, sequence)
    , INDEX command (aggregate_id, event_type, instance_id, sequence) 
        STORING (user_id, service, event_version, payload, aggregate_type, owner)
);