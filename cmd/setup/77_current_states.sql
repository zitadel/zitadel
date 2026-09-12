-- filter_offset stores events2.in_tx_order for sort-key resume, not a row OFFSET.
UPDATE projections.current_states cs
SET filter_offset = e.in_tx_order
FROM eventstore.events2 e
WHERE cs.instance_id = e.instance_id
  AND cs.aggregate_id = e.aggregate_id
  AND cs.aggregate_type = e.aggregate_type
  AND cs."sequence" = e.sequence;
