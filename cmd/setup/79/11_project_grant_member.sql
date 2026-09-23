-- unique_field: projectID:userID:grantID (project.NewAddProjectGrantMemberUniqueConstraint)
-- resource_owner is the project org; also stamp the user's home org and the granted org
WITH page AS MATERIALIZED (
	SELECT m.instance_id, m.project_id, m.grant_id, m.user_id, m.resource_owner, m.user_resource_owner, m.granted_org
	FROM {{.project_grant_members}} m
	WHERE (m.instance_id, m.project_id, m.grant_id, m.user_id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''), COALESCE($1::jsonb->>2, ''), COALESCE($1::jsonb->>3, ''))
	ORDER BY m.instance_id, m.project_id, m.grant_id, m.user_id
	LIMIT $2
),
keys AS (
	SELECT m.instance_id, m.resource_owner, m.user_resource_owner, m.user_id, m.project_id, m.granted_org, m.grant_id,
		m.project_id || ':' || m.user_id || ':' || m.grant_id AS unique_field
	FROM page m
	UNION
	SELECT m.instance_id, m.resource_owner, m.user_resource_owner, m.user_id, m.project_id, m.granted_org, m.grant_id,
		lower(m.project_id || ':' || m.user_id || ':' || m.grant_id)
	FROM page m
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || m.resource_owner, 'org:' || m.user_resource_owner, 'user:' || m.user_id, 'project:' || m.project_id]
		|| CASE WHEN COALESCE(m.granted_org, '') <> '' THEN ARRAY['org:' || m.granted_org] ELSE ARRAY[]::text[] END
		|| CASE WHEN COALESCE(m.grant_id, '') <> '' THEN ARRAY['grant:' || m.grant_id] ELSE ARRAY[]::text[] END
	FROM keys m
	WHERE uc.instance_id = m.instance_id
		AND uc.unique_type = 'project_grant_member'
		AND uc.unique_field = m.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, project_id, grant_id, user_id] FROM page ORDER BY instance_id DESC, project_id DESC, grant_id DESC, user_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
