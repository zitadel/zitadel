-- unique_field: name:resourceOwner (action.NewAddActionNameUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT a.instance_id, a.id, a.resource_owner, a.name
	FROM {{.actions}} a
	WHERE (a.instance_id, a.id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''))
	ORDER BY a.instance_id, a.id
	LIMIT $2
),
keys AS (
	SELECT a.instance_id, a.resource_owner, a.name || ':' || a.resource_owner AS unique_field
	FROM page a
	UNION
	SELECT a.instance_id, a.resource_owner, lower(a.name || ':' || a.resource_owner)
	FROM page a
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.resource_owner]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'action_names'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
