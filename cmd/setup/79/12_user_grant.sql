-- unique_field: resourceOwner:userID:projectID:grantID (usergrant.NewAddUserGrantUniqueConstraint)
-- grant resource_owner is not always the user's org; stamp user, project, and granted orgs
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner, 'org:' || g.resource_owner_user, 'org:' || g.resource_owner_project, 'user:' || g.user_id, 'project:' || g.project_id]
	|| CASE WHEN COALESCE(g.granted_org, '') <> '' THEN ARRAY['org:' || g.granted_org] ELSE ARRAY[]::text[] END
	|| CASE WHEN COALESCE(g.grant_id, '') <> '' THEN ARRAY['grant:' || g.grant_id] ELSE ARRAY[]::text[] END
FROM {{.user_grants}} g
WHERE uc.unique_type = 'user_grant'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id, lower(g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id));
