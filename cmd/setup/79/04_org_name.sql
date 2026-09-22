-- unique_field: org name (org.NewAddOrgNameUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || o.id]
FROM {{.orgs}} o
WHERE uc.unique_type = 'org_name'
	AND uc.instance_id = o.instance_id
	AND uc.owners = '{}'
	AND uc.unique_field IN (o.name, lower(o.name));
