ALTER TABLE IF EXISTS projections.sms_configs3_http ADD COLUMN IF NOT EXISTS signing_key TEXT NULL;
ALTER TABLE IF EXISTS projections.smtp_configs5_http ADD COLUMN IF NOT EXISTS signing_key TEXT NULL;
