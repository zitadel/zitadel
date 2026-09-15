-- unique_field: name:resourceOwner (action.NewAddActionNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || a.resource_owner]
FROM projections.actions3 a
WHERE uc.unique_type = 'action_names'
	AND uc.instance_id = a.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(a.name || ':' || a.resource_owner);
