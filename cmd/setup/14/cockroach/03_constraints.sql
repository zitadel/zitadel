ALTER TABLE eventstore.events2 ALTER COLUMN event_type SET NOT NULL;
ALTER TABLE eventstore.events2 ALTER COLUMN revision SET NOT NULL;
ALTER TABLE eventstore.events2 ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE eventstore.events2 ALTER COLUMN creator SET NOT NULL;
ALTER TABLE eventstore.events2 ALTER COLUMN "owner" SET NOT NULL;
ALTER TABLE eventstore.events2 ALTER COLUMN "position" SET NOT NULL;
ALTER TABLE eventstore.events2 ALTER COLUMN in_tx_order SET NOT NULL;