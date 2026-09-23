-- unique_field: aggregateID:mailTextType:language (policy.NewAddMailTextUniqueConstraint)
WITH page AS MATERIALIZED (
	SELECT t.instance_id, t.aggregate_id, t.type, t.language
	FROM {{.message_texts}} t
	WHERE (t.instance_id, t.aggregate_id, t.type, t.language) > (COALESCE($1::jsonb->>0, ''), COALESCE($1::jsonb->>1, ''), COALESCE($1::jsonb->>2, ''), COALESCE($1::jsonb->>3, ''))
		AND t.aggregate_id <> t.instance_id
	ORDER BY t.instance_id, t.aggregate_id, t.type, t.language
	LIMIT $2
),
keys AS (
	SELECT t.instance_id, t.aggregate_id, t.aggregate_id || ':' || t.type || ':' || t.language AS unique_field
	FROM page t
	UNION
	SELECT t.instance_id, t.aggregate_id, lower(t.aggregate_id || ':' || t.type || ':' || t.language)
	FROM page t
),
upd AS (
	UPDATE eventstore.unique_constraints uc
	SET owners = ARRAY['org:' || k.aggregate_id]
	FROM keys k
	WHERE uc.instance_id = k.instance_id
		AND uc.unique_type = 'mail_text'
		AND uc.unique_field = k.unique_field
		AND uc.owners = '{}'
	RETURNING 1
)
SELECT
	(SELECT ARRAY[instance_id, aggregate_id, type, language] FROM page ORDER BY instance_id DESC, aggregate_id DESC, type DESC, language DESC LIMIT 1),
	(SELECT COUNT(*) FROM upd);
