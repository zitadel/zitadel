BEGIN;
ALTER TABLE eventstore.events2 DROP CONSTRAINT IF EXISTS events2_pkey;
ALTER TABLE eventstore.events2 ADD PRIMARY KEY (instance_id, aggregate_type, aggregate_id, "sequence");
COMMIT;