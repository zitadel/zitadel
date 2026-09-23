-- unique_field: entityID (project.NewAddSAMLConfigEntityIDUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT s.instance_id, s.app_id, s.entity_id
	FROM {{.apps_saml}} s
	WHERE (s.instance_id, s.app_id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY s.instance_id, s.app_id
	LIMIT $2
),
batch AS (
	SELECT s.instance_id, a.id, a.resource_owner, a.project_id, s.entity_id
	FROM page s
	JOIN {{.apps}} a ON a.instance_id = s.instance_id AND a.id = s.app_id
),
keys AS (
	SELECT a.instance_id, a.id, a.resource_owner, a.project_id, a.entity_id AS unique_field
	FROM batch a
	UNION
	SELECT a.instance_id, a.id, a.resource_owner, a.project_id, lower(a.entity_id)
	FROM batch a
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
		AND uc.unique_type = 'entity_ids'
		AND uc.unique_field = a.unique_field
		AND NOT (uc.owners @> ARRAY['app:' || a.id])
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, app_id] FROM page ORDER BY instance_id DESC, app_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
