-- unique_field: org name (org.NewAddOrgNameUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT o.instance_id, o.id, o.name
	FROM {{.orgs}} o
	WHERE (o.instance_id, o.id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''))
	ORDER BY o.instance_id, o.id
	LIMIT $2
),
keys AS (
	SELECT o.instance_id, o.id, o.name AS unique_field
	FROM page o
	UNION
	SELECT o.instance_id, o.id, lower(o.name)
	FROM page o
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.id]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'org_name'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, id] FROM page ORDER BY instance_id DESC, id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
