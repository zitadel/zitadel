-- unique_field: username+resourceOwner (user.NewAddUsernameUniqueConstraint, org-scoped)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || u.resource_owner, 'user:' || u.id]
FROM projections.users14 u
WHERE uc.unique_type = 'usernames'
	AND uc.instance_id = u.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(u.username || u.resource_owner);
