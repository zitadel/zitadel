-- filter_offset stores events2.in_tx_order (reused column; not a row OFFSET).
SELECT
    aggregate_id
    , aggregate_type
    , "sequence"
    , event_date
    , "position"
    , filter_offset
FROM 
    projections.current_states
WHERE
    instance_id = $1
    AND projection_name = $2
FOR NO KEY UPDATE;