-- unique_field: entityID (project.NewAddSAMLConfigEntityIDUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = (
	CASE WHEN uc.owners = '{}'
		THEN ARRAY['org:' || a.resource_owner, 'project:' || a.project_id]
		ELSE uc.owners
	END
) || ARRAY['app:' || a.id]
FROM {{.apps}} a
JOIN {{.apps_saml}} s ON s.app_id = a.id AND s.instance_id = a.instance_id
WHERE uc.unique_type = 'entity_ids'
	AND uc.instance_id = a.instance_id
	AND NOT (uc.owners @> ARRAY['app:' || a.id])
	AND uc.unique_field IN (s.entity_id, lower(s.entity_id));
