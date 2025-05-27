SELECT data_type
FROM information_schema.columns
WHERE table_schema = 'eventstore'
AND table_name = 'events2'
AND column_name = 'in_tx_order';
