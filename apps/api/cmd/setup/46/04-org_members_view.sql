CREATE OR REPLACE VIEW eventstore.org_members AS
SELECT instance_id, aggregate_id as org_id, object_id as user_id, text_value as role
FROM eventstore.fields
WHERE aggregate_type = 'org'
AND object_type = 'org_member_role'
AND field_name = 'org_role';
