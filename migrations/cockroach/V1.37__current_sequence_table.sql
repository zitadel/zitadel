DELETE FROM management.current_sequences WHERE aggregate_type <> '';
DELETE FROM auth.current_sequences WHERE aggregate_type <> '';
DELETE FROM authz.current_sequences WHERE aggregate_type <> '';
DELETE FROM adminapi.current_sequences WHERE aggregate_type <> '';
DELETE FROM notification.current_sequences WHERE aggregate_type <> '';

BEGIN;

ALTER TABLE management.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE auth.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE authz.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE adminapi.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE notification.current_sequences DROP CONSTRAINT "primary";

ALTER TABLE management.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name);
ALTER TABLE auth.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name);
ALTER TABLE authz.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name);
ALTER TABLE adminapi.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name);
ALTER TABLE notification.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name);

COMMIT;

BEGIN;

ALTER TABLE management.current_sequences DROP COLUMN aggregate_type;
ALTER TABLE auth.current_sequences DROP COLUMN aggregate_type;
ALTER TABLE authz.current_sequences DROP COLUMN aggregate_type;
ALTER TABLE adminapi.current_sequences DROP COLUMN aggregate_type;
ALTER TABLE notification.current_sequences DROP COLUMN aggregate_type;

COMMIT;