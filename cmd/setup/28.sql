-- TODO: whats the correct name of this table?
CREATE TABLE eventstore.lookups (
    instance_id TEXT NOT NULL
    , resource_owner TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , field_name TEXT NOT NULL
    , "value" JSONB
    , text_value TEXT
    , CONSTRAINT one_of_values CHECK (("value" IS NULL) <> (text_value IS NULL))
);