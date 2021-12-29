ALTER TABLE zitadel.projections.users ADD COLUMN type INT2;

-- human is 1
-- machine is 2
WITH doa AS (
    SELECT u.id, IF(h.user_id IS NULL, 2, 1) as type
    FROM zitadel.projections.users u 
    LEFT JOIN zitadel.projections.users_humans h
        ON h.user_id = u.id
    LEFT JOIN zitadel.projections.users_machines m
        ON m.user_id = u.id 
)
UPDATE zitadel.projections.users SET type = doa.type FROM doa WHERE doa.id = zitadel.projections.users.id;
