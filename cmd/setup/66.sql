-- Create table to track per-table DML event times for service ping telemetry
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE TABLE IF NOT EXISTS analytics.service_ping_resource_events
(
    id SERIAL PRIMARY KEY,
    instance_id TEXT NOT NULL,
    table_name TEXT NOT NULL,
    parent_type TEXT NOT NULL,
    parent_id TEXT NOT NULL,
    event TEXT NOT NULL, -- 'INSERT' | 'UPDATE' | 'DELETE'
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- TODO - the purpose of this is really just to figure out if onboarding works. We don't necessarily want
-- 50,000 events when 50,000 users are created. WE need to determine if we should handle this via scheduled truncation,
-- only adding to a certain number of rows, or only logging during onboarding somehow.

-- Trigger function to insert an entry whenever a row is inserted/updated/deleted
-- Arguments (TG_ARGV):
-- 1. parent_type (TEXT)
-- 2. instance_id column name (TEXT)
-- 3. parent_id column name (TEXT)
CREATE OR REPLACE FUNCTION analytics.record_service_ping_resource_event()
    RETURNS trigger
    LANGUAGE 'plpgsql' VOLATILE
AS $$
DECLARE
    tg_table_name TEXT := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
    tg_parent_type TEXT := TG_ARGV[0];
    tg_instance_id_column TEXT := TG_ARGV[1];
    tg_parent_id_column TEXT := TG_ARGV[2];

    tg_instance_id TEXT;
    tg_parent_id TEXT;

    select_ids TEXT := format('SELECT ($1).%I, ($1).%I', tg_instance_id_column, tg_parent_id_column);
BEGIN
    IF (TG_OP = 'INSERT') THEN
        EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING NEW;
        INSERT INTO analytics.service_ping_resource_events(instance_id, table_name, parent_type, parent_id, event)
        VALUES (tg_instance_id, tg_table_name, tg_parent_type, tg_parent_id, 'INSERT');
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING NEW;
        INSERT INTO analytics.service_ping_resource_events(instance_id, table_name, parent_type, parent_id, event)
        VALUES (tg_instance_id, tg_table_name, tg_parent_type, tg_parent_id, 'UPDATE');
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING OLD;
        INSERT INTO analytics.service_ping_resource_events(instance_id, table_name, parent_type, parent_id, event)
        VALUES (tg_instance_id, tg_table_name, tg_parent_type, tg_parent_id, 'DELETE');
        RETURN OLD;
    END IF;
END
$$;


