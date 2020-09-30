package scripts

const V16OriginAllowList = `
BEGIN;

ALTER TABLE management.applications ADD COLUMN origin_allow_list TEXT ARRAY;
ALTER TABLE auth.applications ADD COLUMN origin_allow_list TEXT ARRAY;
ALTER TABLE authz.applications ADD COLUMN origin_allow_list TEXT ARRAY;

TRUNCATE TABLE management.applications;
TRUNCATE TABLE auth.applications;
TRUNCATE TABLE authz.applications;

UPDATE management.current_sequences set current_sequence = 0 where view_name = 'management.applications';
UPDATE auth.current_sequences set current_sequence = 0 where view_name = 'auth.applications';
UPDATE authz.current_sequences set current_sequence = 0 where view_name = 'authz.applications';

COMMIT;
`
