-- unique_field: name+resourceOwner (idpconfig.NewAddIDPConfigNameUniqueConstraint, templates)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || t.resource_owner, 'idp:' || t.id]
FROM projections.idp_templates6 t
WHERE uc.unique_type = 'idp_config_names'
	AND uc.instance_id = t.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(t.name || t.resource_owner);
