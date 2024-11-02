select
    id,
    creation_date,
    login_client,
    client_id,
    scope,
    redirect_uri,
    prompt,
    ui_locales,
    login_hint,
    max_age,
    hint_user_id
from projections.auth_requests %s
where id = $1 and instance_id = $2
limit 1;
