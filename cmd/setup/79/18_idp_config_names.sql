-- unique_field: name+resourceOwner (idpconfig.NewAddIDPConfigNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || i.resource_owner, 'idp:' || i.id]
FROM {{.idps}} i
WHERE uc.unique_type = 'idp_config_names'
	AND uc.instance_id = i.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (i.name || i.resource_owner, lower(i.name || i.resource_owner));
