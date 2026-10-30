-- filter_offset is events2.in_tx_order when setup 77 is Done, else a row OFFSET.
SELECT
    aggregate_id
    , aggregate_type
    , "sequence"
    , event_date
    , "position"
    , filter_offset
    , in_tx_order
FROM 
    projections.current_states
WHERE
    instance_id = $1
    AND projection_name = $2
FOR NO KEY UPDATE;
