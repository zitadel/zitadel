-- unique_field: name+resourceOwner (project.NewAddProjectNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || p.resource_owner, 'project:' || p.id]
FROM {{.projects}} p
WHERE uc.unique_type = 'project_names'
	AND uc.instance_id = p.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (p.name || p.resource_owner, lower(p.name || p.resource_owner));
