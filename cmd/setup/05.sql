CREATE INDEX IF NOT EXISTS current_sequences_instance_id_idx ON adminapi.current_sequences (instance_id);
CREATE INDEX IF NOT EXISTS current_sequences_instance_id_idx ON auth.current_sequences (instance_id);
CREATE INDEX IF NOT EXISTS current_sequences_instance_id_idx ON projections.current_sequences (instance_id);

CREATE INDEX IF NOT EXISTS failed_events_instance_id_idx ON adminapi.failed_events (instance_id);
CREATE INDEX IF NOT EXISTS failed_events_instance_id_idx ON auth.failed_events (instance_id);
CREATE INDEX IF NOT EXISTS failed_events_instance_id_idx ON projections.failed_events (instance_id);

ALTER TABLE adminapi.failed_events ADD COLUMN IF NOT EXISTS last_failed TIMESTAMPTZ;
ALTER TABLE auth.failed_events ADD COLUMN IF NOT EXISTS last_failed TIMESTAMPTZ;
ALTER TABLE projections.failed_events ADD COLUMN IF NOT EXISTS last_failed TIMESTAMPTZ;
