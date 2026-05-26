ALTER TABLE IF EXISTS projections.lockout_policies3 
ADD COLUMN IF NOT EXISTS show_remaining_lockout_time BOOLEAN NOT NULL DEFAULT FALSE;