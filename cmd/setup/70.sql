ALTER TABLE IF EXISTS projections.password_complexity_policies2 ADD COLUMN IF NOT EXISTS history_count INT8 NOT NULL DEFAULT 0;
