-- unique_field: roleKey:projectID (project.NewAddProjectRoleUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT r.instance_id, r.project_id, r.role_key, r.resource_owner
	FROM {{.project_roles}} r
	WHERE (r.instance_id, r.project_id, r.role_key) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''), COALESCE(($1::text[])[3], ''))
	ORDER BY r.instance_id, r.project_id, r.role_key
	LIMIT $2
),
keys AS (
	SELECT r.instance_id, r.resource_owner, r.project_id, r.role_key || ':' || r.project_id AS unique_field
	FROM page r
	UNION
	SELECT r.instance_id, r.resource_owner, r.project_id, lower(r.role_key || ':' || r.project_id)
	FROM page r
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.resource_owner, 'project:' || k.project_id]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'project_role'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, project_id, role_key] FROM page ORDER BY instance_id DESC, project_id DESC, role_key DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
