-- Create table to track per-table DML event times for service ping telemetry
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE TABLE IF NOT EXISTS analytics.events
(
    id BIGSERIAL PRIMARY KEY,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    event_data JSONB NOT NULL
);