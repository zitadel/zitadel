ALTER TABLE IF EXISTS projections.lockout_policies3 
ADD COLUMN IF NOT EXISTS auto_unlock_after_min BIGINT NOT NULL DEFAULT 0;