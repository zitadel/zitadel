CREATE DATABASE IF NOT EXISTS zitadel;

CREATE SCHEMA IF NOT EXISTS zitadel.eventstore_v5;

CREATE TABLE IF NOT EXISTS zitadel.eventstore_v5.events (
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

    , CONSTRAINT agg_seq_uniq UNIQUE(instance_id, aggregate_id, sequence desc)
    , INDEX command (aggregate_id, event_type, instance_id, sequence desc) 
        STORING (user_id, service, event_version, payload, aggregate_type, owner)
);