package scripts

const V14Compliance = `
BEGIN;

ALTER TABLE management.applications ADD COLUMN oidc_version SMALLINT;
ALTER TABLE management.applications ADD COLUMN none_compliant BOOLEAN;
ALTER TABLE management.applications ADD COLUMN compliance_problems TEXT ARRAY;
ALTER TABLE management.applications ADD COLUMN dev_mode BOOLEAN;

ALTER TABLE auth.applications ADD COLUMN oidc_version SMALLINT;
ALTER TABLE auth.applications ADD COLUMN none_compliant BOOLEAN;
ALTER TABLE auth.applications ADD COLUMN compliance_problems TEXT ARRAY;
ALTER TABLE auth.applications ADD COLUMN dev_mode BOOLEAN;

ALTER TABLE authz.applications ADD COLUMN oidc_version SMALLINT;
ALTER TABLE authz.applications ADD COLUMN none_compliant BOOLEAN;
ALTER TABLE authz.applications ADD COLUMN compliance_problems TEXT ARRAY;
ALTER TABLE authz.applications ADD COLUMN dev_mode BOOLEAN;

COMMIT;
`
