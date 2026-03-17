-- represents an event to be created.
DO $$ BEGIN
    CREATE TYPE eventstore.command2 AS (
        instance_id TEXT
        , aggregate_type TEXT
        , aggregate_id TEXT
        , command_type TEXT
        , revision INT2
        , payload JSONB
        , creator TEXT
        , owner TEXT
        , written_by_relational BOOLEAN
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE IF EXISTS eventstore.events2 ADD COLUMN IF NOT EXISTS written_by_relational BOOLEAN NOT NULL DEFAULT false;
