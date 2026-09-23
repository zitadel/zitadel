-- unique_field: resourceOwner:userID:projectID:grantID (usergrant.NewAddUserGrantUniqueConstraint)
-- grant resource_owner is not always the user's org; stamp user, project, and granted orgs
WITH page AS MATERIALIZED (
	SELECT g.instance_id, g.id, g.resource_owner, g.resource_owner_user, g.resource_owner_project, g.user_id, g.project_id, g.granted_org, g.grant_id
	FROM {{.user_grants}} g
	WHERE (g.instance_id, g.id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''))
	ORDER BY g.instance_id, g.id
	LIMIT $2
),
keys AS (
	SELECT g.instance_id, g.resource_owner, g.resource_owner_user, g.resource_owner_project, g.user_id, g.project_id, g.granted_org, g.grant_id,
		g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id AS unique_field
	FROM page g
	UNION
	SELECT g.instance_id, g.resource_owner, g.resource_owner_user, g.resource_owner_project, g.user_id, g.project_id, g.granted_org, g.grant_id,
		lower(g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id)
	FROM page g
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || g.resource_owner, 'org:' || g.resource_owner_user, 'org:' || g.resource_owner_project, 'user:' || g.user_id, 'project:' || g.project_id]
		|| CASE WHEN COALESCE(g.granted_org, '') <> '' THEN ARRAY['org:' || g.granted_org] ELSE ARRAY[]::text[] END
		|| CASE WHEN COALESCE(g.grant_id, '') <> '' THEN ARRAY['grant:' || g.grant_id] ELSE ARRAY[]::text[] END
	FROM keys g
	WHERE uc.instance_id = g.instance_id
		AND uc.unique_type = 'user_grant'
		AND uc.unique_field = g.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
