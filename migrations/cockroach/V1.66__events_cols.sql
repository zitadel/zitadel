BEGIN;

ALTER TABLE eventstore.events
    RENAME COLUMN previous_sequence TO previous_aggregate_sequence,
    ADD COLUMN previous_aggregate_type_sequence INT8, 
    ADD CONSTRAINT prev_agg_type_seq_unique UNIQUE(previous_aggregate_type_sequence);

COMMIT;