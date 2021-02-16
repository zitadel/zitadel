BEGIN;

TRUNCATE TABLE authz.user_memberships;

INSERT INTO authz.user_memberships (
    user_id,
    member_type,
    aggregate_id,
    object_id,
    roles,
    display_name,
    resource_owner,
    resource_owner_name,
    creation_date,
    change_date,
    sequence
)
SELECT
    user_id,
    member_type,
    aggregate_id,
    object_id,
    roles,
    display_name,
    resource_owner,
    resource_owner_name,
    creation_date,
    change_date,
    sequence
FROM auth.user_memberships;

INSERT INTO authz.current_sequences (
     view_name,
     event_timestamp,
     current_sequence,
     last_successful_spooler_run,
     aggregate_type
     )
SELECT 'authz.user_memberships',
    event_timestamp,
    current_sequence,
    last_successful_spooler_run,
    aggregate_type
FROM auth.current_sequences
WHERE view_name = 'auth.user_memberships';

COMMIT;