-- project_grant_members lists the roles users hold as members of project grants.
-- The project grant member events store them on the project aggregate,
-- with the user ID as object ID and the grant ID as field name.
CREATE OR REPLACE VIEW eventstore.project_grant_members AS
SELECT instance_id, aggregate_id AS project_id, field_name AS grant_id, object_id AS user_id, text_value AS role
FROM eventstore.fields
WHERE aggregate_type = 'project'
AND object_type = 'project_grant_member_role';
