-- Create table to track per-table DML event times for service ping telemetry
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE TABLE IF NOT EXISTS analytics.events
(
    id BIGSERIAL PRIMARY KEY,
    instance_id TEXT NOT NULL,
    event JSONB NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);