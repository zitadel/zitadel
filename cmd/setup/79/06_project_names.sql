-- unique_field: name+resourceOwner (project.NewAddProjectNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || p.resource_owner, 'project:' || p.id]
FROM projections.projects4 p
WHERE uc.unique_type = 'project_names'
	AND uc.instance_id = p.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(p.name || p.resource_owner);
