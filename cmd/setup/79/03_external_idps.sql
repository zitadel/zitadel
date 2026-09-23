-- unique_field: idpConfigID+externalUserID (user.NewAddUserIDPLinkUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT l.instance_id, l.idp_id, l.external_user_id, l.resource_owner, l.user_id
	FROM {{.idp_user_links}} l
	WHERE (l.instance_id, l.idp_id, l.external_user_id) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''), COALESCE($1::jsonb->>2, ''))
	ORDER BY l.instance_id, l.idp_id, l.external_user_id
	LIMIT $2
),
keys AS (
	SELECT l.instance_id, l.resource_owner, l.idp_id, l.user_id, l.idp_id || l.external_user_id AS unique_field
	FROM page l
	UNION
	SELECT l.instance_id, l.resource_owner, l.idp_id, l.user_id, lower(l.idp_id || l.external_user_id)
	FROM page l
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || l.resource_owner, 'idp:' || l.idp_id, 'user:' || l.user_id]
	FROM keys l
	WHERE uc.instance_id = l.instance_id
		AND uc.unique_type = 'external_idps'
		AND uc.unique_field = l.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, idp_id, external_user_id] FROM page ORDER BY instance_id DESC, idp_id DESC, external_user_id DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
