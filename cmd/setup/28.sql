ALTER TABLE IF EXISTS projections.sms_configs2_twilio ADD COLUMN IF NOT EXISTS verify_service_sid TEXT not null;
