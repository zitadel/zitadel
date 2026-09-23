-- unique_field: username (user.NewAddUsernameUniqueConstraint, instance-scoped)
WITH page AS MATERIALIZED (
	SELECT u.instance_id, u.id, u.resource_owner, u.username
	FROM {{.users}} u
	WHERE (u.instance_id, u.id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY u.instance_id, u.id
	LIMIT $2
),
batch AS (
	SELECT u.instance_id, u.id, u.resource_owner, u.username
	FROM page u
	LEFT JOIN {{.domain_policies}} org_pol
		ON org_pol.instance_id = u.instance_id AND org_pol.id = u.resource_owner
	LEFT JOIN {{.domain_policies}} inst_pol
		ON inst_pol.instance_id = u.instance_id AND inst_pol.is_default
	WHERE NOT COALESCE(org_pol.user_login_must_be_domain, inst_pol.user_login_must_be_domain, false)
),
keys AS (
	SELECT u.instance_id, u.id, u.resource_owner, u.username AS unique_field
	FROM batch u
	UNION
	SELECT u.instance_id, u.id, u.resource_owner, lower(u.username)
	FROM batch u
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || u.resource_owner, 'user:' || u.id]
	FROM keys u
	WHERE uc.instance_id = u.instance_id
		AND uc.unique_type = 'usernames'
		AND uc.unique_field = u.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
