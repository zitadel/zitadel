-- represents an event to be created.
DO $$ BEGIN
    CREATE TYPE eventstore.command AS (
        instance_id TEXT
        , aggregate_type TEXT
        , aggregate_id TEXT
        , command_type TEXT
        , revision INT2
        , payload JSONB
        , creator TEXT
        , owner TEXT
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
