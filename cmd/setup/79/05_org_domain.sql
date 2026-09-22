-- unique_field: domain (org.NewAddOrgDomainUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || d.org_id]
FROM {{.org_domains}} d
WHERE uc.unique_type = 'org_domain'
	AND uc.instance_id = d.instance_id
	AND uc.owners = '{}'
	AND d.is_verified
	AND uc.unique_field IN (d.domain, lower(d.domain));
