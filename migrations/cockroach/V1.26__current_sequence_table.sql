ALTER TABLE management.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE auth.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE authz.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE adminapi.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';
ALTER TABLE notification.current_sequences ADD COLUMN aggregate_type STRING NOT NULL DEFAULT '';

ALTER TABLE management.current_sequences ALTER PRIMARY KEY USING COLUMNS (view_name, aggregate_type);
ALTER TABLE auth.current_sequences ALTER PRIMARY KEY USING COLUMNS (view_name, aggregate_type);
ALTER TABLE authz.current_sequences ALTER PRIMARY KEY USING COLUMNS (view_name, aggregate_type);
ALTER TABLE adminapi.current_sequences ALTER PRIMARY KEY USING COLUMNS (view_name, aggregate_type);
ALTER TABLE notification.current_sequences ALTER PRIMARY KEY USING COLUMNS (view_name, aggregate_type);