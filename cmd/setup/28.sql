-- TODO: whats the correct name of this table?
CREATE TABLE eventstore.fields (
    instance_id TEXT NOT NULL
    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL
    , field_name TEXT NOT NULL
    , "value" JSONB
    , text_value TEXT
    , CHECK (("value" IS NULL) <> (text_value IS NULL))
);