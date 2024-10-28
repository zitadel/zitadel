select s.id, s.creation_date, s.change_date, s.resource_owner, s.sequence, s.session_id, s.user_id, s.client_id,
       c.back_channel_logout_uri
from projections.notification_oidc_sessions s
left join projections.apps7_oidc_configs c on s.instance_id = c.instance_id and s.client_id = c.client_id
where s.instance_id = $1
and s.session_id = $2
