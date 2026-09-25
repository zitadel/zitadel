-- unique_field: domain (org.NewAddOrgDomainUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT d.org_id, d.domain, d.instance_id
	FROM {{.org_domains}} d
	WHERE (d.org_id, d.domain, d.instance_id) > (COALESCE(($1::text[])[1], ''), COALESCE(($1::text[])[2], ''), COALESCE(($1::text[])[3], ''))
		AND d.is_verified
	ORDER BY d.org_id, d.domain, d.instance_id
	LIMIT $2
),
keys AS (
	SELECT d.instance_id, d.org_id, d.domain AS unique_field
	FROM page d
	UNION
	SELECT d.instance_id, d.org_id, lower(d.domain)
	FROM page d
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.org_id]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'org_domain'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[org_id, domain, instance_id] FROM page ORDER BY org_id DESC, domain DESC, instance_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
