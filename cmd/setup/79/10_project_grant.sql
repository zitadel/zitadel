-- unique_field: grantedOrgID:projectID (project.NewAddProjectGrantUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner, 'org:' || g.granted_org_id, 'project:' || g.project_id]
FROM projections.project_grants4 g
WHERE uc.unique_type = 'project_grant'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(g.granted_org_id || ':' || g.project_id);
