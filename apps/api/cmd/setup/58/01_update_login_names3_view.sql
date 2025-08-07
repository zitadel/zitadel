CREATE OR REPLACE VIEW projections.login_names3 AS
    SELECT
        u.id AS user_id
        , CASE
            WHEN p.must_be_domain THEN CONCAT(u.user_name, '@', d.name)
            ELSE u.user_name
        END AS login_name
        , COALESCE(d.is_primary, TRUE) AS is_primary
        , u.instance_id
    FROM
        projections.login_names3_users AS u
    LEFT JOIN LATERAL (
        SELECT
            must_be_domain
            , is_default
        FROM
            projections.login_names3_policies AS p
        WHERE
            (
                p.instance_id = u.instance_id
                AND NOT p.is_default
                AND p.resource_owner = u.resource_owner
            ) OR (
                p.instance_id = u.instance_id
                AND p.is_default
            )
        ORDER BY
            p.is_default -- custom first
        LIMIT 1
    ) AS p ON TRUE
    LEFT JOIN 
        projections.login_names3_domains d
        ON 
            p.must_be_domain 
            AND u.resource_owner = d.resource_owner 
            AND u.instance_id = d.instance_id
