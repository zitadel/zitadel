with auth_methods as (select user_id, method_type, token_id, state, instance_id, name
                      from projections.user_auth_methods4
                      where instance_id = $1
                        and user_id = $2
                      LIMIT 1),
     verified_auth_methods as (select method_type from auth_methods where state = 2)
select u.id,
       u.creation_date,
       u.change_date,
       u.resource_owner,
       u.state                                                                        as user_state,
       au.password_set,
       h.password_change_required,
       au.password_change,
       au.last_login,
       u.username                                                                     as user_name,
       (SELECT array_agg(ll.login_name) login_names
        FROM projections.login_names3 ll
        WHERE u.id = ll.user_id
          and u.instance_id = ll.instance_id
        GROUP BY ll.user_id, ll.instance_id)
                                                                                      as login_names,
       l.login_name,
       h.first_name,
       h.last_name,
       h.nick_name,
       h.display_name,
       h.preferred_language,
       h.gender,
       h.email,
       h.is_email_verified,
       h.phone,
       h.is_phone_verified,
       o.country,
       o.locality,
       o.postal_code,
       o.region,
       o.street_address,
       (select coalesce((select state from auth_methods where method_type = 1), 0))   as otp_state,
       case
           when exists (select true from verified_auth_methods where method_type = 3) then 2
           when exists (select true from verified_auth_methods where method_type = 2) then 1
           else 0 end                                                                 as mfa_max_set_up,
       au.mfa_init_skipped,
       u.sequence,
       au.init_required,
       au.username_change_required,
       m.name                                                                         as machine_name,
       m.description                                                                  as machine_description,
       u.type                                                                         as user_type,
       (select JSONB_AGG(json_build_object('webAuthNTokenId', token_id, 'webAuthNTokenName', name, 'state', state))
        from auth_methods
        where method_type = 2)                                                        as u2f_tokens,
       (select JSONB_AGG(json_build_object('webAuthNTokenId', token_id, 'webAuthNTokenName', name, 'state', state))
        from auth_methods
        where method_type = 3)                                                        as passwordless_tokens,
       h.avatar_key,
       au.passwordless_init_required,
       au.password_init_required,
       u.instance_id,
       o.owner_removed,
       (select exists (select true from verified_auth_methods where method_type = 6)) as otp_sms_added,
       (select exists (select true from verified_auth_methods where method_type = 7)) as otp_email_added
from projections.users12 u
         left join projections.users12_humans h on u.instance_id = h.instance_id and u.id = h.user_id
         left join projections.login_names3 l
                   on l.instance_id = u.instance_id and l.user_id = u.id and l.is_primary = true
         left join projections.users12_machines m on m.instance_id =
                                                     u.instance_id and m.user_id = u.id
         left join auth.users3 au on au.instance_id = u.instance_id and au.id = u.id
         left join auth.users2 o on o.id = u.id and o.instance_id = u.instance_id
where u.instance_id = $1
  and u.id = $2
LIMIT 1;