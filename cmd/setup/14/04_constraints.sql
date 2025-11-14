ALTER TABLE eventstore.events2 ALTER COLUMN event_type SET NOT NULL,
                               ALTER COLUMN revision SET NOT NULL,
                               ALTER COLUMN created_at SET NOT NULL,
                               ALTER COLUMN creator SET NOT NULL,
                               ALTER COLUMN "owner" SET NOT NULL,
                               ALTER COLUMN "position" SET NOT NULL,
                               ALTER COLUMN in_tx_order SET NOT NULL;