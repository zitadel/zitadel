CREATE INDEX CONCURRENTLY IF NOT EXISTS events2_current_sequence2
    ON eventstore.events2 USING btree
    (aggregate_id ASC, aggregate_type ASC, instance_id ASC, sequence DESC);
