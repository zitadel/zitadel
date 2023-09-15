ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS "position" xid8;
ALTER TABLE eventstore.events ALTER COLUMN "position" SET NOT NULL,
                              ALTER COLUMN "position" TYPE xid8 USING pg_current_xact_id();

ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS in_tx_order INTEGER;
ALTER TABLE eventstore.events ALTER COLUMN in_tx_order SET NOT NULL,
                              ALTER COLUMN in_tx_order TYPE INTEGER USING event_sequence;