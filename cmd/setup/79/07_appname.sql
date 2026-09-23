-- unique_field: name:projectID (project.NewAddApplicationUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = (
	CASE WHEN uc.owners = '{}'
		THEN ARRAY['org:' || a.resource_owner, 'project:' || a.project_id]
		ELSE uc.owners
	END
) || ARRAY['app:' || a.id]
FROM {{.apps}} a
WHERE uc.unique_type = 'appname'
	AND uc.instance_id = a.instance_id
	AND NOT (uc.owners @> ARRAY['app:' || a.id])
	AND uc.unique_field IN (a.name || ':' || a.project_id, lower(a.name || ':' || a.project_id));
