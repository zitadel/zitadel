ALTER TABLE eventstore.events ALTER COLUMN "position" SET NOT NULL,
                              ALTER COLUMN "position" TYPE DECIMAL USING COALESCE("position", EXTRACT(EPOCH FROM creation_date));