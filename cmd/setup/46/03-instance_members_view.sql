CREATE OR REPLACE VIEW eventstore.instance_members AS
SELECT instance_id, object_id as user_id, text_value as role
FROM eventstore.fields
WHERE aggregate_type = 'instance'
AND object_type = 'instance_member_role'
AND field_name = 'instance_role';
