BEGIN;

ALTER TABLE eventstore.events
    RENAME COLUMN previous_sequence TO previous_aggregate_sequence,
    ADD COLUMN previous_aggregate_type_sequence INT8, 
    ADD CONSTRAINT prev_agg_type_seq_unique UNIQUE(previous_aggregate_type_sequence);

COMMIT;

SET CLUSTER SETTING kv.closed_timestamp.target_duration = '2m';

BEGIN;
WITH data AS (
    SELECT 
    event_sequence, 
    LAG(event_sequence) 
        OVER (
            PARTITION BY aggregate_type 
            ORDER BY event_sequence
        ) as prev_seq,
    aggregate_type
    FROM eventstore.events
    ORDER BY event_sequence
) UPDATE eventstore.events 
    SET previous_aggregate_type_sequence = data.prev_seq
    FROM data 
    WHERE data.event_sequence = events.event_sequence;
COMMIT;

SET CLUSTER SETTING kv.closed_timestamp.target_duration TO DEFAULT;

-- validation by hand:
-- SELECT 
--     event_sequence, 
--     previous_aggregate_sequence, 
--     previous_aggregate_type_sequence,
--     aggregate_type,
--     aggregate_id,
--     event_type
-- FROM eventstore.events ORDER BY event_sequence;