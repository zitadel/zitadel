-- unique_field: username (user.NewAddUsernameUniqueConstraint, instance-scoped)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || u.resource_owner, 'user:' || u.id]
FROM {{.users}} u
LEFT JOIN {{.domain_policies}} org_pol
	ON org_pol.instance_id = u.instance_id AND org_pol.id = u.resource_owner
LEFT JOIN {{.domain_policies}} inst_pol
	ON inst_pol.instance_id = u.instance_id AND inst_pol.is_default
WHERE uc.unique_type = 'usernames'
	AND uc.instance_id = u.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (u.username, lower(u.username))
	AND NOT COALESCE(org_pol.user_login_must_be_domain, inst_pol.user_login_must_be_domain, false);
