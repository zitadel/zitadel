WITH found_users AS (
  SELECT DISTINCT
    u.id
    , u.instance_id
    , u.resource_owner
    , u.user_name
    , COALESCE(p_custom.must_be_domain, p_default.must_be_domain) as must_be_domain
  FROM 
    projections.login_names3_users u
    LEFT JOIN projections.login_names3_policies p_custom
      ON  u.instance_id = p_custom.instance_id
        AND p_custom.instance_id = $4 AND p_custom.resource_owner = u.resource_owner
    LEFT JOIN projections.login_names3_policies p_default
      ON  u.instance_id = p_default.instance_id
      AND p_default.instance_id = $4 AND p_default.is_default IS TRUE
      AND (
        (COALESCE(p_custom.must_be_domain, p_default.must_be_domain) IS TRUE AND u.user_name_lower = $1)
         OR (COALESCE(p_custom.must_be_domain, p_default.must_be_domain) IS FALSE AND u.user_name_lower = $3)
      )
  JOIN 
    projections.login_names3_domains d
    ON 
      u.instance_id = d.instance_id
      AND u.resource_owner = d.resource_owner
      AND CASE WHEN COALESCE(p_custom.must_be_domain, p_default.must_be_domain) THEN d.name_lower = $2 ELSE TRUE END
  WHERE 
    u.instance_id = $4
    AND u.user_name_lower IN (
      $1, 
      $3
    )
),
login_names AS (SELECT 
  fu.id user_id
  , fu.instance_id
  , fu.resource_owner
  , fu.user_name
  , d.name domain_name
  , d.is_primary
  , fu.must_be_domain
  , CASE WHEN fu.must_be_domain
      THEN concat(fu.user_name, '@', d.name)
      ELSE fu.user_name
    END login_name
  FROM 
    found_users fu
  JOIN 
    projections.login_names3_domains d
    ON 
      fu.instance_id = d.instance_id
      AND fu.resource_owner = d.resource_owner
)
SELECT 
  u.id
  , u.creation_date
  , u.change_date
  , u.resource_owner
  , u.sequence
  , u.state
  , u.type
  , u.username
  , (SELECT array_agg(ln.login_name)::TEXT[] login_names FROM login_names ln WHERE fu.id = ln.user_id GROUP BY ln.user_id, ln.instance_id) login_names
  , (SELECT ln.login_name login_names_lower FROM login_names ln WHERE fu.id = ln.user_id AND ln.is_primary IS TRUE) preferred_login_name
  , h.user_id
  , h.first_name
  , h.last_name
  , h.nick_name
  , h.display_name
  , h.preferred_language
  , h.gender
  , h.avatar_key
  , h.email
  , h.is_email_verified
  , h.phone
  , h.is_phone_verified
  , h.password_change_required
  , m.user_id
  , m.name
  , m.description
  , m.secret
  , m.access_token_type
  , count(*) OVER ()
FROM found_users fu
JOIN
  projections.users12 u
  ON
    fu.id = u.id
    AND fu.instance_id = u.instance_id
LEFT JOIN
  projections.users12_humans h
  ON
    fu.id = h.user_id
    AND fu.instance_id = h.instance_id
LEFT JOIN
  projections.users12_machines m
  ON
    fu.id = m.user_id
    AND fu.instance_id = m.instance_id
WHERE 
  u.instance_id = $4
;