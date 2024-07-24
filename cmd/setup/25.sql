ALTER TABLE IF EXISTS projections.users13_notifications ADD COLUMN IF NOT EXISTS verified_email_lower TEXT GENERATED ALWAYS AS (lower(verified_email)) STORED;
CREATE INDEX IF NOT EXISTS users13_notifications_email_search ON projections.users13_notifications (instance_id, verified_email_lower);
