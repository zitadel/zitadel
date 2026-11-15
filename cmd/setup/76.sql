ALTER TABLE eventstore.unique_constraints
    ADD COLUMN IF NOT EXISTS claimant_type TEXT,
    ADD COLUMN IF NOT EXISTS claimant_id TEXT;
