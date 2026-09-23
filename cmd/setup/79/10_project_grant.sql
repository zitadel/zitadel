-- unique_field: grantedOrgID:projectID (project.NewAddProjectGrantUniqueConstraint)
-- resource_owner is the project org and granted_org_id is the granted org; they differ
WITH page AS MATERIALIZED (
	SELECT g.instance_id, g.grant_id, g.resource_owner, g.granted_org_id, g.project_id
	FROM {{.project_grants}} g
	WHERE (g.instance_id, g.grant_id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''))
	ORDER BY g.instance_id, g.grant_id
	LIMIT $2
),
keys AS (
	SELECT g.instance_id, g.resource_owner, g.granted_org_id, g.project_id, g.grant_id, g.granted_org_id || ':' || g.project_id AS unique_field
	FROM page g
	UNION
	SELECT g.instance_id, g.resource_owner, g.granted_org_id, g.project_id, g.grant_id, lower(g.granted_org_id || ':' || g.project_id)
	FROM page g
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || g.resource_owner, 'org:' || g.granted_org_id, 'project:' || g.project_id, 'grant:' || g.grant_id]
	FROM keys g
	WHERE uc.instance_id = g.instance_id
		AND uc.unique_type = 'project_grant'
		AND uc.unique_field = g.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, grant_id] FROM page ORDER BY instance_id DESC, grant_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
