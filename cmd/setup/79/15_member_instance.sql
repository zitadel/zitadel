-- unique_field: instanceID:userID (member.NewAddMemberUniqueConstraint, instance aggregate)
-- instance members live on the instance aggregate, so also stamp the user's home org
WITH page AS MATERIALIZED (
	SELECT m.instance_id, m.id, m.user_id, m.user_resource_owner
	FROM {{.instance_members}} m
	WHERE (m.instance_id, m.id, m.user_id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''), COALESCE(($1::text[])[3], ''))
	ORDER BY m.instance_id, m.id, m.user_id
	LIMIT $2
),
keys AS (
	SELECT m.instance_id, m.user_id, m.user_resource_owner, m.id || ':' || m.user_id AS unique_field
	FROM page m
	UNION
	SELECT m.instance_id, m.user_id, m.user_resource_owner, lower(m.id || ':' || m.user_id)
	FROM page m
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['user:' || k.user_id, 'org:' || k.user_resource_owner]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'member'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id, user_id] FROM page ORDER BY instance_id DESC, id DESC, user_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
