-- is used to set the value for previous data
SET enable_experimental_alter_column_type_general = true;
SET enable_implicit_transaction_for_batch_statements = off;
SET sql_safe_updates = false;

ALTER TABLE eventstore.events ALTER COLUMN in_tx_order TYPE INTEGER USING COALESCE(in_tx_order, event_sequence);

RESET enable_experimental_alter_column_type_general;
RESET enable_implicit_transaction_for_batch_statements;
RESET sql_safe_updates;