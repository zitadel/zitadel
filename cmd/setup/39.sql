DELETE FROM eventstore.fields
WHERE aggregate_type = 'org'
AND aggregate_id IN (
	SELECT aggregate_id
	FROM eventstore.events2
	WHERE instance_id = $1
	AND aggregate_type = 'org'
	AND event_type = 'org.removed'
);
