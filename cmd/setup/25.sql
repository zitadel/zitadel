ALTER TABLE IF EXISTS projections.users10_notifications ADD COLUMN IF NOT EXISTS verified_email_lower TEXT GENERATED ALWAYS AS (lower(verified_email)) STORED;
CREATE INDEX CONCURRENTLY IF NOT EXISTS users10_notifications_email_search ON projections.users10_notifications (instance_id, verified_email_lower);
