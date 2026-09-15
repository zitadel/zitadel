-- unique_field: domain (org.NewAddOrgDomainUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || d.org_id]
FROM projections.org_domains2 d
WHERE uc.unique_type = 'org_domain'
	AND uc.instance_id = d.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(d.domain);
