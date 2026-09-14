-- unique_field: instanceID:userID (member.NewAddMemberUniqueConstraint, instance aggregate)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['user:' || m.user_id]
FROM projections.instance_members4 m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.id || ':' || m.user_id);
