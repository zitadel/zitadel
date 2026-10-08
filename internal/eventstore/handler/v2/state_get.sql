-- filter_offset is either events2.in_tx_order (setup 77 already Done) or a
-- row OFFSET count. in_tx_order is the events2 ordinal; the BEFORE UPDATE
-- trigger clears it unless this statement opted in via UPDATE OF in_tx_order.
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
