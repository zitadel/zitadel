ALTER TABLE management.applications ADD COLUMN clock_skew BIGINT;
ALTER TABLE auth.applications ADD COLUMN clock_skew BIGINT;
ALTER TABLE authz.applications ADD COLUMN clock_skew BIGINT;