ALTER TABLE IF EXISTS projections.instances
    ADD COLUMN IF NOT EXISTS audit_log_retention INTERVAL,
    ADD COLUMN IF NOT EXISTS block BOOLEAN;
;

