ALTER TABLE eventstore.events ALTER COLUMN in_tx_order SET NOT NULL,
                              ALTER COLUMN in_tx_order TYPE INTEGER USING COALESCE(in_tx_order, event_sequence);