-- unique_field: name+resourceOwner (idpconfig.NewAddIDPConfigNameUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT i.instance_id, i.id, i.resource_owner, i.name
	FROM {{.idps}} i
	WHERE (i.instance_id, i.id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY i.instance_id, i.id
	LIMIT $2
),
keys AS (
	SELECT i.instance_id, i.resource_owner, i.id, i.name || i.resource_owner AS unique_field
	FROM page i
	UNION
	SELECT i.instance_id, i.resource_owner, i.id, lower(i.name || i.resource_owner)
	FROM page i
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || i.resource_owner, 'idp:' || i.id]
	FROM keys i
	WHERE uc.instance_id = i.instance_id
		AND uc.unique_type = 'idp_config_names'
		AND uc.unique_field = i.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
