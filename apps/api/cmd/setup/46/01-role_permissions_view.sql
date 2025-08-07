CREATE OR REPLACE VIEW eventstore.role_permissions AS
SELECT instance_id, aggregate_id, object_id as role, text_value as permission
FROM eventstore.fields
WHERE aggregate_type = 'permission'
AND object_type = 'role_permission'
AND field_name = 'permission';
