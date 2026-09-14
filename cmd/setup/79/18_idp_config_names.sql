-- unique_field: name+resourceOwner (idpconfig.NewAddIDPConfigNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || i.resource_owner, 'idp:' || i.id]
FROM projections.idps3 i
WHERE uc.unique_type = 'idp_config_names'
	AND uc.instance_id = i.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(i.name || i.resource_owner);
