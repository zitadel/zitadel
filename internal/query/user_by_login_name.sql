with login_names as (
  select 
    u.id user_id
    , u.instance_id
    , u.resource_owner
    , u.user_name
    , d.name domain_name
    , d.is_primary
    , p.must_be_domain
    , CASE WHEN p.must_be_domain 
        THEN concat(u.user_name, '@', d.name)
        ELSE u.user_name
      END login_name
  from 
    projections.login_names2_users u
  join 
    projections.login_names2_domains d
    on 
      (u.user_name = $1 OR u.user_name = $3)
      AND u.instance_id = $4
      AND u.instance_id = d.instance_id
      AND u.resource_owner = d.resource_owner
  join
    projections.login_names2_policies p
    on
      u.instance_id = p.instance_id
      -- AND u.resource_owner = p.resource_owner
      AND (
        (p.is_default is TRUE AND p.instance_id = $4)
        OR (p.instance_id = $4 AND p.resource_owner = u.resource_owner)
      )
      AND (
          (p.must_be_domain is true and u.user_name = $1 and d.name = $2)
          or (p.must_be_domain is false and u.user_name = $3)
      )
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
  , (select array_agg(ln.login_name)::TEXT[] login_names from login_names ln group by ln.user_id, ln.instance_id) loginnames
  , (select ln.login_name login_names_lower from login_names ln where ln.is_primary is true) preferred_login_name
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
  , m.user_id
  , m.name
  , m.description
  , m.has_secret
  , m.access_token_type
  , count(*) OVER ()
FROM login_names ln
JOIN
  projections.users8 u
  ON
    ln.user_id = u.id
    AND ln.instance_id = u.instance_id
LEFT JOIN
  projections.users8_humans h
  ON
    ln.user_id = h.user_id
    AND ln.instance_id = h.instance_id
LEFT JOIN
  projections.users8_machines m
  ON
    ln.user_id = m.user_id
    AND ln.instance_id = m.instance_id
WHERE 
  u.instance_id = $2
LIMIT 1
;