BEGIN;

TRUNCATE auth.login_policies;
TRUNCATE management.login_policies;

UPDATE auth.current_sequences set current_sequence = 0 where view_name = 'auth.login_policies';
UPDATE management.current_sequences set current_sequence = 0 where view_name = 'management.login_policies';

COMMIT;