-- unique_field: aggregateID:mailTextType:language (policy.NewAddMailTextUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || t.aggregate_id]
FROM {{.message_texts}} t
WHERE uc.unique_type = 'mail_text'
	AND uc.instance_id = t.instance_id
	AND uc.owners = '{}'
	AND t.aggregate_id <> t.instance_id
	AND uc.unique_field IN (t.aggregate_id || ':' || t.type || ':' || t.language, lower(t.aggregate_id || ':' || t.type || ':' || t.language));
