ALTER TABLE management.applications ADD COLUMN origin_allow_list TEXT ARRAY;
ALTER TABLE auth.applications ADD COLUMN origin_allow_list TEXT ARRAY;
ALTER TABLE authz.applications ADD COLUMN origin_allow_list TEXT ARRAY;

DELETE FROM management.applications;
DELETE FROM auth.applications;
DELETE FROM authz.applications;

UPDATE management.current_sequences set current_sequence = 0 where view_name = 'management.applications';
UPDATE auth.current_sequences set current_sequence = 0 where view_name = 'auth.applications';
UPDATE authz.current_sequences set current_sequence = 0 where view_name = 'authz.applications';
