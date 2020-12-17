ALTER TABLE management.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE auth.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE authz.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE adminapi.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE notification.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';

BEGIN;

ALTER TABLE management.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE auth.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE authz.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE adminapi.current_sequences DROP CONSTRAINT "primary";
ALTER TABLE notification.current_sequences DROP CONSTRAINT "primary";

ALTER TABLE management.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name, aggregate_type);
ALTER TABLE auth.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name, aggregate_type);
ALTER TABLE authz.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name, aggregate_type);
ALTER TABLE adminapi.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name, aggregate_type);
ALTER TABLE notification.current_sequences ADD CONSTRAINT "primary" PRIMARY KEY (view_name, aggregate_type);

COMMIT;