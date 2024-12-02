DELETE FROM eventstore.fields
WHERE aggregate_type = 'org'
AND aggregate_id NOT IN (
	SELECT id
	FROM projections.orgs1
);
