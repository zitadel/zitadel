CREATE OR REPLACE VIEW eventstore.project_members AS
SELECT instance_id, aggregate_id as project_id, object_id as user_id, text_value as role
FROM eventstore.fields
WHERE aggregate_type = 'project'
AND object_type = 'project_member_role'
AND field_name = 'project_role';
