-- unique_field: projectID:userID (member.NewAddMemberUniqueConstraint, project aggregate)
-- stamp the project org and the user's home org (they differ for cross-org members)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || m.resource_owner, 'org:' || m.user_resource_owner, 'user:' || m.user_id, 'project:' || m.project_id]
FROM projections.project_members4 m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.project_id || ':' || m.user_id);
