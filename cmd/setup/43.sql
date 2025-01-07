CREATE INDEX CONCURRENTLY events2_current_sequence2
    ON eventstore.events2 USING btree
    (aggregate_id ASC, aggregate_type ASC, instance_id ASC, sequence DESC);

DROP INDEX eventstore.events2_current_sequence;