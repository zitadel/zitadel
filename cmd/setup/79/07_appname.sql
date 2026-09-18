-- unique_field: name:projectID (project.NewAddApplicationUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || a.resource_owner, 'project:' || a.project_id]
FROM {{.apps}} a
WHERE uc.unique_type = 'appname'
	AND uc.instance_id = a.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (a.name || ':' || a.project_id, lower(a.name || ':' || a.project_id));
