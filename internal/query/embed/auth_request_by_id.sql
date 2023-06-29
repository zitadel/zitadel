select
    id,
    creation_date,
    client_id,
    scope,
    redirect_uri,
    prompt,
    ui_locales,
    login_hint,
    max_age,
    hint_user_id
from projections.auth_requests
where id = $1
limit 1;
