ALTER TABLE adminapi.failed_events ALTER PRIMARY KEY USING COLUMNS (instance_id, view_name, failed_sequence);
ALTER TABLE auth.failed_events ALTER PRIMARY KEY USING COLUMNS (instance_id, view_name, failed_sequence);
ALTER TABLE projections.failed_events ALTER PRIMARY KEY USING COLUMNS (instance_id, projection_name, failed_sequence);

ALTER TABLE adminapi.failed_events ADD COLUMN last_failed TIMESTAMPTZ;
ALTER TABLE auth.failed_events ADD COLUMN last_failed TIMESTAMPTZ;
ALTER TABLE projections.failed_events ADD COLUMN last_failed TIMESTAMPTZ;

ALTER TABLE adminapi.current_sequences ALTER PRIMARY KEY USING COLUMNS (instance_id, view_name);
ALTER TABLE auth.current_sequences ALTER PRIMARY KEY USING COLUMNS (instance_id, view_name);
ALTER TABLE projections.current_sequences ALTER PRIMARY KEY USING COLUMNS (instance_id, projection_name, aggregate_type);
