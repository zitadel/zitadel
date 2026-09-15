-- unique_field: idpConfigID+externalUserID (user.NewAddUserIDPLinkUniqueConstraint)
UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || l.resource_owner, 'idp:' || l.idp_id, 'user:' || l.user_id]
FROM projections.idp_user_links3 l
WHERE uc.unique_type = 'external_idps'
	AND uc.instance_id = l.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(l.idp_id || l.external_user_id);
