-- unique_field: name+resourceOwner (idpconfig.NewAddIDPConfigNameUniqueConstraint, templates)
WITH page AS MATERIALIZED (
	SELECT t.instance_id, t.id, t.resource_owner, t.name
	FROM {{.idp_templates}} t
	WHERE (t.instance_id, t.id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY t.instance_id, t.id
	LIMIT $2
),
keys AS (
	SELECT t.instance_id, t.resource_owner, t.id, t.name || t.resource_owner AS unique_field
	FROM page t
	UNION
	SELECT t.instance_id, t.resource_owner, t.id, lower(t.name || t.resource_owner)
	FROM page t
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || t.resource_owner, 'idp:' || t.id]
	FROM keys t
	WHERE uc.instance_id = t.instance_id
		AND uc.unique_type = 'idp_config_names'
		AND uc.unique_field = t.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
