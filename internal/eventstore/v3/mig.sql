CREATE DATABASE IF NOT EXISTS zitadel;

CREATE SCHEMA IF NOT EXISTS zitadel.eventstore_v3;

CREATE TABLE IF NOT EXISTS zitadel.eventstore_v3.events (
    -- aggregate information
    aggregate_id STRING NOT NULL
    , aggregate_type STRING NOT NULL
    , owner STRING NOT NULL
    , instance_id STRING NOT NULL
    -- editor metadata
    , user_id STRING NOT NULL
    , service STRING NOT NULL
    -- event metadata
    , event_type STRING NOT NULL
    , event_version SMALLINT NOT NULL
    , creation_date TIMESTAMPTZ NOT NULL
    , payload JSONB
    , region STRING

    , PRIMARY KEY (instance_id, aggregate_id, creation_date DESC)
    , INDEX command (aggregate_id, aggregate_type, event_type, instance_id, creation_date) 
        STORING (owner, user_id, service, event_version, payload)
);