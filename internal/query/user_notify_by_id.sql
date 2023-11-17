with login_names AS (select 
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
    u.id = $1
    AND u.instance_id = $2
    AND u.instance_id = d.instance_id
    AND u.resource_owner = d.resource_owner
  join lateral (
    SELECT * FROM projections.login_names2_policies p
    where
      u.instance_id = p.instance_id
      AND (
        (p.is_default = TRUE AND p.instance_id = $2)
        OR (p.instance_id = $2 AND p.resource_owner = u.resource_owner)
      )
      ORDER BY is_default
      LIMIT 1
  ) p on true
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
  , (select array_agg(ln.login_name)::TEXT[] login_names from login_names ln group by ln.user_id, ln.instance_id) login_names
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
  , n.user_id
  , n.last_email
  , n.verified_email
  , n.last_phone
  , n.verified_phone
  , n.password_set
  , count(*) OVER ()
FROM projections.users8 u
LEFT JOIN
  projections.users8_humans h
  ON
    u.id = h.user_id
    AND u.instance_id = h.instance_id
LEFT JOIN
  projections.users8_notifications m
  ON
    u.id = m.user_id
    AND u.instance_id = m.instance_id
WHERE 
  u.id = $1
  AND u.instance_id = $2
LIMIT 1
;