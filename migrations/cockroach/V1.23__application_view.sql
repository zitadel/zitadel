ALTER TABLE management.applications ADD COLUMN id_token_userinfo_assertion BOOLEAN;
ALTER TABLE auth.applications ADD COLUMN id_token_userinfo_assertion BOOLEAN;
ALTER TABLE authz.applications ADD COLUMN id_token_userinfo_assertion BOOLEAN;

ALTER TABLE management.applications ADD COLUMN clock_skew BIGINT;
ALTER TABLE auth.applications ADD COLUMN clock_skew BIGINT;
ALTER TABLE authz.applications ADD COLUMN clock_skew BIGINT;
