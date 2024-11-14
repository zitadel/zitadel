-- index is used for filtering for the current sequence of the aggregate
CREATE INDEX CONCURRENTLY IF NOT EXISTS e_push_idx ON eventstore.events2(instance_id, aggregate_type, aggregate_id, owner, sequence DESC)
