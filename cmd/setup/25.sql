ALTER TABLE IF EXISTS projections.users12_notifications ADD COLUMN IF NOT EXISTS verified_email_lower TEXT GENERATED ALWAYS AS (lower(verified_email)) STORED;
CREATE INDEX IF NOT EXISTS users12_notifications_email_search ON projections.users12_notifications (instance_id, verified_email_lower);
