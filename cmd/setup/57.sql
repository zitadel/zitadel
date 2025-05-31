ALTER TABLE IF EXISTS projections.sessions8
ADD COLUMN IF NOT EXISTS mfa_recovery_code_checked_at TIMESTAMPTZ; 