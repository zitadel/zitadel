-- unique_field: name:resourceOwner (group.NewAddGroupNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner]
FROM projections.groups1 g
WHERE uc.unique_type = 'group_name'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(g.name || ':' || g.resource_owner);
