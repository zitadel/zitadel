-- unique_field: orgID:userID (member.NewAddMemberUniqueConstraint, org aggregate)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || m.org_id, 'user:' || m.user_id]
FROM projections.org_members4 m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.org_id || ':' || m.user_id);
