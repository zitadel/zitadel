-- unique_field: grantedOrgID:projectID (project.NewAddProjectGrantUniqueConstraint)
-- resource_owner is the project org and granted_org_id is the granted org; they differ
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner, 'org:' || g.granted_org_id, 'project:' || g.project_id, 'grant:' || g.grant_id]
FROM {{.project_grants}} g
WHERE uc.unique_type = 'project_grant'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (g.granted_org_id || ':' || g.project_id, lower(g.granted_org_id || ':' || g.project_id));
