ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS "position" DECIMAL;
ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS in_tx_order INTEGER;
