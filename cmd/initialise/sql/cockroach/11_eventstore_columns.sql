-- is used to set the value for previous data
SET enable_experimental_alter_column_type_general = true;
SET enable_implicit_transaction_for_batch_statements = off;
SET sql_safe_updates = false;

ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS "position" DECIMAL;
ALTER TABLE eventstore.events ALTER COLUMN "position" TYPE DECIMAL USING creation_date::DECIMAL;
ALTER TABLE eventstore.events ALTER COLUMN "position" SET NOT NULL;

ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS in_tx_order INTEGER;
ALTER TABLE eventstore.events ALTER COLUMN in_tx_order TYPE INTEGER USING event_sequence;
ALTER TABLE eventstore.events ALTER COLUMN in_tx_order SET NOT NULL;