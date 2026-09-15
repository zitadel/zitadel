-- unique_field: roleKey:projectID (project.NewAddProjectRoleUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || r.resource_owner, 'project:' || r.project_id]
FROM projections.project_roles4 r
WHERE uc.unique_type = 'project_role'
	AND uc.instance_id = r.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(r.role_key || ':' || r.project_id);
