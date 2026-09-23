-- unique_field: name:projectID (project.NewAddApplicationUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT a.instance_id, a.id, a.resource_owner, a.project_id, a.name
	FROM {{.apps}} a
	WHERE (a.instance_id, a.id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''))
	ORDER BY a.instance_id, a.id
	LIMIT $2
),
keys AS (
	SELECT a.instance_id, a.id, a.resource_owner, a.project_id, a.name || ':' || a.project_id AS unique_field
	FROM page a
	UNION
	SELECT a.instance_id, a.id, a.resource_owner, a.project_id, lower(a.name || ':' || a.project_id)
	FROM page a
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = (
		CASE WHEN uc.owners = '{}'
			THEN ARRAY['org:' || a.resource_owner, 'project:' || a.project_id]
			ELSE uc.owners
		END
	) || ARRAY['app:' || a.id]
	FROM keys a
	WHERE uc.instance_id = a.instance_id
		AND uc.unique_type = 'appname'
		AND uc.unique_field = a.unique_field
		AND NOT (uc.owners @> ARRAY['app:' || a.id])
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
