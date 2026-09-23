-- unique_field: projectID:userID (member.NewAddMemberUniqueConstraint, project aggregate)
-- stamp the project org and the user's home org (they differ for cross-org members)
WITH page AS MATERIALIZED (
	SELECT m.instance_id, m.project_id, m.user_id, m.resource_owner, m.user_resource_owner
	FROM {{.project_members}} m
	WHERE (m.instance_id, m.project_id, m.user_id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''), COALESCE($1::jsonb->>2, ''))
	ORDER BY m.instance_id, m.project_id, m.user_id
	LIMIT $2
),
keys AS (
	SELECT m.instance_id, m.resource_owner, m.user_resource_owner, m.user_id, m.project_id, m.project_id || ':' || m.user_id AS unique_field
	FROM page m
	UNION
	SELECT m.instance_id, m.resource_owner, m.user_resource_owner, m.user_id, m.project_id, lower(m.project_id || ':' || m.user_id)
	FROM page m
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.resource_owner, 'org:' || k.user_resource_owner, 'user:' || k.user_id, 'project:' || k.project_id]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'member'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, project_id, user_id] FROM page ORDER BY instance_id DESC, project_id DESC, user_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
