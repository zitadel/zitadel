-- unique_field: name+resourceOwner (idpconfig.NewAddIDPConfigNameUniqueConstraint, templates)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || t.resource_owner, 'idp:' || t.id]
FROM {{.idp_templates}} t
WHERE uc.unique_type = 'idp_config_names'
	AND uc.instance_id = t.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (t.name || t.resource_owner, lower(t.name || t.resource_owner));
