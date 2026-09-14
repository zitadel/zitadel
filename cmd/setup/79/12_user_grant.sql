-- unique_field: resourceOwner:userID:projectID:grantID (usergrant.NewAddUserGrantUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner, 'user:' || g.user_id, 'project:' || g.project_id] || CASE WHEN COALESCE(g.grant_id, '') <> '' THEN ARRAY['grant:' || g.grant_id] ELSE ARRAY[]::text[] END
FROM projections.user_grants5 g
WHERE uc.unique_type = 'user_grant'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id);
