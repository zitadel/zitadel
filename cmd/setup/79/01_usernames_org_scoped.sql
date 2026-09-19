-- unique_field: username+resourceOwner (user.NewAddUsernameUniqueConstraint, org-scoped)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || u.resource_owner, 'user:' || u.id]
FROM {{.users}} u
WHERE uc.unique_type = 'usernames'
	AND uc.instance_id = u.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (u.username || u.resource_owner, lower(u.username || u.resource_owner));
