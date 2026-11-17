ALTER TABLE eventstore.unique_constraints ADD COLUMN IF NOT EXISTS owners TEXT[] NOT NULL DEFAULT '{}';
