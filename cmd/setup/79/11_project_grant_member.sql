-- unique_field: projectID:userID:grantID (project.NewAddProjectGrantMemberUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || m.resource_owner, 'user:' || m.user_id, 'project:' || m.project_id] || CASE WHEN COALESCE(m.grant_id, '') <> '' THEN ARRAY['grant:' || m.grant_id] ELSE ARRAY[]::text[] END
FROM projections.project_grant_members4 m
WHERE uc.unique_type = 'project_grant_member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.project_id || ':' || m.user_id || ':' || m.grant_id);
