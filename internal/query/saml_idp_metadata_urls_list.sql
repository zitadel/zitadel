SELECT s.id,
    t.instance_id,
    t.resource_owner,
    t.owner_type,
    s.metadata_url
FROM projections.idp_templates6_saml s
JOIN projections.idp_templates6 t
    ON s.id = t.id
    AND s.instance_id = t.instance_id
WHERE t.type = $1
    AND t.state = $2
    AND s.metadata_url IS NOT NULL
    AND s.metadata_url <> ''
    AND (t.instance_id, s.id) > ($3, $4)
ORDER BY t.instance_id, s.id
LIMIT $5;
