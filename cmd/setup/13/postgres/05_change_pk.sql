ALTER TABLE eventstore.events DROP CONSTRAINT IF EXISTS events_pkey;
ALTER TABLE eventstore.events ADD PRIMARY KEY (instance_id, aggregate_type, aggregate_id, event_sequence);