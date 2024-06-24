CREATE TABLE eventstore.search (
    instance_id TEXT NOT NULL
    , resource_owner TEXT NOT NULL

    , aggregate_type TEXT NOT NULL
    , aggregate_id TEXT NOT NULL

    , object_type TEXT NOT NULL
    , object_id TEXT NOT NULL
    , object_revision INT4 NOT NULL -- we use INT4 here because PSQL does not support unsigned numbers
    
    , field_name TEXT NOT NULL
    , number_value NUMERIC
    , text_value TEXT
    
    , CONSTRAINT one_of_values CHECK (num_nonnulls(number_value, text_value) = 1)
);

CREATE INDEX IF NOT EXISTS search_number_value_idx ON eventstore.search (instance_id, object_type, object_revision, field_name, number_value) INCLUDE (resource_owner, aggregate_type, aggregate_id, object_id);
CREATE INDEX IF NOT EXISTS search_text_value_idx ON eventstore.search (instance_id, object_type, object_revision, field_name, text_value) INCLUDE (resource_owner, aggregate_type, aggregate_id, object_id);
CREATE INDEX IF NOT EXISTS search_object_idx ON eventstore.search (instance_id, object_type, object_id, object_revision) INCLUDE (resource_owner, aggregate_type, aggregate_id);
