-- unique_field: name+resourceOwner (project.NewAddProjectNameUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT p.instance_id, p.id, p.resource_owner, p.name
	FROM {{.projects}} p
	WHERE (p.instance_id, p.id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY p.instance_id, p.id
	LIMIT $2
),
keys AS (
	SELECT p.instance_id, p.id, p.resource_owner, p.name || p.resource_owner AS unique_field
	FROM page p
	UNION
	SELECT p.instance_id, p.id, p.resource_owner, lower(p.name || p.resource_owner)
	FROM page p
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.resource_owner, 'project:' || k.id]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'project_names'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
