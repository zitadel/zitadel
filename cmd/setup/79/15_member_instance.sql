-- unique_field: instanceID:userID (member.NewAddMemberUniqueConstraint, instance aggregate)
-- instance members live on the instance aggregate, so also stamp the user's home org
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['user:' || m.user_id, 'org:' || m.user_resource_owner]
FROM {{.instance_members}} m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (m.id || ':' || m.user_id, lower(m.id || ':' || m.user_id));
