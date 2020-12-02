ALTER TABLE management.current_sequences ADD COLUMN last_successful_spooler_run TIMESTAMPTZ;
ALTER TABLE auth.current_sequences ADD COLUMN last_successful_spooler_run TIMESTAMPTZ;
ALTER TABLE authz.current_sequences ADD COLUMN last_successful_spooler_run TIMESTAMPTZ;
ALTER TABLE adminapi.current_sequences ADD COLUMN last_successful_spooler_run TIMESTAMPTZ;
ALTER TABLE notification.current_sequences ADD COLUMN last_successful_spooler_run TIMESTAMPTZ;

ALTER TABLE management.current_sequences RENAME COLUMN timestamp TO event_timestamp;
ALTER TABLE auth.current_sequences RENAME COLUMN timestamp TO event_timestamp;
ALTER TABLE authz.current_sequences RENAME COLUMN timestamp TO event_timestamp;
ALTER TABLE adminapi.current_sequences RENAME COLUMN timestamp TO event_timestamp;
ALTER TABLE notification.current_sequences RENAME COLUMN timestamp TO event_timestamp;
