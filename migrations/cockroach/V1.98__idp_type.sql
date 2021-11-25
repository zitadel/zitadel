ALTER TABLE zitadel.projections.idps ADD COLUMN type INT2;

-- jwt-type is 2
-- oidc-type is 0
WITH doa AS (
    SELECT i.id, IF(o.idp_id IS NULL, 0, 2) as type
    FROM projections.idps i 
    LEFT JOIN projections.idps_oidc_config o 
        ON o.idp_id = i.id 
    LEFT JOIN projections.idps_jwt_config j 
        ON j.idp_id = i.id
)
UPDATE zitadel.projections.silvan_idps SET type = doa.type FROM doa WHERE doa.id = zitadel.projections.silvan_idps.id;
