-- unique_field: username+resourceOwner (user.NewAddUsernameUniqueConstraint, org-scoped)
WITH page AS MATERIALIZED (
	SELECT u.instance_id, u.id, u.resource_owner, u.username
	FROM {{.users}} u
	WHERE (u.instance_id, u.id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''))
	ORDER BY u.instance_id, u.id
	LIMIT $2
),
keys AS (
	SELECT u.instance_id, u.id, u.resource_owner, u.username || u.resource_owner AS unique_field
	FROM page u
	UNION
	SELECT u.instance_id, u.id, u.resource_owner, lower(u.username || u.resource_owner)
	FROM page u
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
