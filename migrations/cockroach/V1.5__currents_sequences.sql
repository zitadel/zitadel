BEGIN;

ALTER TABLE authz.current_sequences ADD COLUMN timestamp TIMESTAMPTZ;
ALTER TABLE auth.current_sequences ADD COLUMN timestamp TIMESTAMPTZ;
ALTER TABLE management.current_sequences ADD COLUMN timestamp TIMESTAMPTZ;
ALTER TABLE notification.current_sequences ADD COLUMN timestamp TIMESTAMPTZ;
ALTER TABLE adminapi.current_sequences ADD COLUMN timestamp TIMESTAMPTZ;

COMMIT;