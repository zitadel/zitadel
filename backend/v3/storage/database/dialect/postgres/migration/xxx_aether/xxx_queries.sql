-- get inherited settings of type 'notification' for a specific collection
WITH RECURSIVE collection_hierarchy AS (
    SELECT instance_id, id AS collection_id, parent_id, 0 AS level
        FROM aether.collections
        WHERE instance_id = 'instance1'
            AND id = 'org1_1_1_1'
    
    UNION ALL
    
    SELECT c.instance_id, c.id AS collection_id, c.parent_id, level + 1
        FROM aether.collections c
        INNER JOIN collection_hierarchy ch ON c.instance_id = ch.instance_id AND c.id = ch.parent_id
), properties AS (
    SELECT 
        s.collection_id
        , sp.*
        , CASE
            WHEN s.collection_id IS NULL THEN 
                ch.level + 1
            ELSE
                ch.level
            END AS level
    FROM collection_hierarchy ch
    JOIN aether.settings s
        ON s.instance_id = ch.instance_id 
        AND (s.collection_id = ch.collection_id OR (ch.parent_id IS NULL AND s.collection_id IS NULL))
    LEFT JOIN aether.setting_properties sp 
        ON sp.instance_id = s.instance_id 
        AND sp.setting_id = s.id
    WHERE s.type = 'notification'
    ORDER BY ch.level DESC, s.collection_id NULLS FIRST
) 
SELECT DISTINCT ON (path) 
    path
    , value
    , setting_id
FROM properties
ORDER BY path, level
;


-- search users
WITH found_users AS (
    SELECT DISTINCT up.instance_id, up.user_id
    FROM aether.user_properties up
    WHERE up.instance_id = 'instance1'
        AND (
            (
                (up.path = 'username' AND up.value @> to_jsonb('user1name'::TEXT))
                OR (up.path = 'profile.email' AND up.value @> to_jsonb('user1@example.com'::TEXT))
            )
            OR
            (
                up.path = 'id' AND up.value @> to_jsonb(3::NUMERIC)
            )
        )
)
SELECT u.* 
FROM aether.users u
JOIN found_users fu ON u.instance_id = fu.instance_id AND u.id = fu.user_id
;