CREATE INDEX CONCURRENTLY IF NOT EXISTS f_aggregate_object_type_idx ON eventstore.fields (aggregate_type, aggregate_id, object_type);
