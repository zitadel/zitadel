-- is used to set the value for previous data
SET enable_experimental_alter_column_type_general = true;
SET enable_implicit_transaction_for_batch_statements = off;
SET sql_safe_updates = false;

ALTER TABLE eventstore.events ALTER COLUMN "position" TYPE DECIMAL USING COALESCE("position", creation_date::DECIMAL);

RESET enable_experimental_alter_column_type_general;
RESET enable_implicit_transaction_for_batch_statements;
RESET sql_safe_updates;