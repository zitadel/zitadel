UPDATE projections.current_states
SET filter_offset = in_tx_order
WHERE in_tx_order IS NOT NULL;
