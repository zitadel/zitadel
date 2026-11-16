-- unique_field: name:resourceOwner (group.NewAddGroupNameUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT g.instance_id, g.id, g.resource_owner, g.name
	FROM {{.groups}} g
	WHERE (g.instance_id, g.id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY g.instance_id, g.id
	LIMIT $2
),
keys AS (
	SELECT g.instance_id, g.resource_owner, g.name || ':' || g.resource_owner AS unique_field
	FROM page g
	UNION
	SELECT g.instance_id, g.resource_owner, lower(g.name || ':' || g.resource_owner)
	FROM page g
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.resource_owner]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'group_name'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
