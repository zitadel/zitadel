package scripts

const V12Views = `
CREATE TABLE management.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE management.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (view_name)
);

CREATE TABLE management.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE management.projects (
    project_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    project_name TEXT,
    project_state SMALLINT,
    resource_owner TEXT,
    sequence BIGINT,

    PRIMARY KEY (project_id)
);

CREATE TABLE management.project_grants (
    grant_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    project_id TEXT,
    project_name TEXT,
    org_name TEXT,
    project_state SMALLINT,
    resource_owner TEXT,
    org_id TEXT,
    granted_role_keys TEXT Array,
    sequence BIGINT,
    resource_owner_name TEXT,

    PRIMARY KEY (grant_id)
);

CREATE TABLE management.project_roles (
    project_id TEXT,
    role_key TEXT,
    display_name TEXT,
    resource_owner TEXT,
    org_id TEXT,
    group_name TEXT,

    creation_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (org_id, project_id, role_key)
);

CREATE TABLE management.project_members (
    user_id TEXT,
    project_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    user_name TEXT,
    email_address TEXT,
    first_name TEXT,
    last_name TEXT,
    roles TEXT ARRAY,
    display_name TEXT,
    sequence BIGINT,

    PRIMARY KEY (project_id, user_id)
);

CREATE TABLE management.project_grant_members (
    user_id TEXT,
    grant_id TEXT,
    project_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    user_name TEXT,
    email_address TEXT,
    first_name TEXT,
    last_name TEXT,
    roles TEXT ARRAY,
    display_name TEXT,
    sequence BIGINT,

    PRIMARY KEY (grant_id, user_id)
);

CREATE TABLE management.applications (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    app_state SMALLINT,
    resource_owner TEXT,
    app_name TEXT,
    project_id TEXT,
    app_type SMALLINT,
    is_oidc BOOLEAN,
    oidc_client_id TEXT,
    oidc_redirect_uris TEXT ARRAY,
    oidc_response_types SMALLINT ARRAY,
    oidc_grant_types SMALLINT ARRAY,
    oidc_application_type SMALLINT,
    oidc_auth_method_type SMALLINT,
    oidc_post_logout_redirect_uris TEXT ARRAY,

    PRIMARY KEY (id)
);

CREATE TABLE management.users (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    user_state SMALLINT,
    last_login TIMESTAMPTZ,
    password_change TIMESTAMPTZ,
    user_name TEXT,
    login_names TEXT ARRAY,
    preferred_login_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nick_Name TEXT,
    display_name TEXT,
    preferred_language TEXT,
    gender SMALLINT,
    email TEXT,
    is_email_verified BOOLEAN,
    phone TEXT,
    is_phone_verified BOOLEAN,
    country TEXT,
    locality TEXT,
    postal_code TEXT,
    region TEXT,
    street_address TEXT,
    otp_state SMALLINT,
    sequence BIGINT,
    password_set BOOLEAN,
    password_change_required BOOLEAN,
    mfa_max_set_up SMALLINT,
    mfa_init_skipped TIMESTAMPTZ,
    init_required BOOLEAN,

    PRIMARY KEY (id)
);

CREATE TABLE management.user_grants (
    id TEXT,
    resource_owner TEXT,
    project_id TEXT,
    user_id TEXT,
    org_name TEXT,
    project_name TEXT,
    user_name TEXT,
    display_name TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT,
    role_keys TEXT Array,
    grant_id TEXT,

    grant_state SMALLINT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

CREATE TABLE management.org_domains (
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    domain TEXT,
    org_id TEXT,
    verified BOOLEAN,
    primary_domain BOOLEAN,

    PRIMARY KEY (org_id, domain)
);

CREATE TABLE auth.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE auth.current_sequences (
    view_name TEXT,
    timestamp TIMESTAMPTZ,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE auth.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE auth.auth_requests (
    id TEXT,
    request JSONB,
    code TEXT,
    request_type smallint,

    PRIMARY KEY (id)
);

CREATE TABLE auth.users (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    user_state SMALLINT,
    password_set BOOLEAN,
    password_change_required BOOLEAN,
    password_change TIMESTAMPTZ,
    last_login TIMESTAMPTZ,
    user_name TEXT,
    login_names TEXT ARRAY,
    preferred_login_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nick_name TEXT,
    display_name TEXT,
    preferred_language TEXT,
    gender SMALLINT,
    email TEXT,
    is_email_verified BOOLEAN,
    phone TEXT,
    is_phone_verified BOOLEAN,
    country TEXT,
    locality TEXT,
    postal_code TEXT,
    region TEXT,
    street_address TEXT,
    otp_state SMALLINT,
    mfa_max_set_up SMALLINT,
    mfa_init_skipped TIMESTAMPTZ,
    sequence BIGINT,
    init_required BOOLEAN,

    PRIMARY KEY (id)
);

CREATE TABLE auth.user_sessions (
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    state SMALLINT,
    user_agent_id TEXT,
    user_id TEXT,
    user_name TEXT,
    password_verification TIMESTAMPTZ,
    mfa_software_verification TIMESTAMPTZ,
    mfa_hardware_verification TIMESTAMPTZ,
    sequence BIGINT,
    mfa_software_verification_type SMALLINT,
    mfa_hardware_verification_type SMALLINT,
    user_display_name TEXT,
    login_name TEXT,

    PRIMARY KEY (user_agent_id, user_id)
);

CREATE TABLE auth.tokens (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    application_id TEXT,
    user_agent_id TEXT,
    user_id TEXT,
    expiration TIMESTAMPTZ,
    sequence BIGINT,
    scopes TEXT ARRAY,
    audience TEXT ARRAY,

    PRIMARY KEY (id)
);


CREATE TABLE notification.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE notification.current_sequences (
    view_name TEXT,
    timestamp TIMESTAMPTZ,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE notification.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE notification.notify_users (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    nick_Name TEXT,
    display_name TEXT,
    preferred_language TEXT,
    gender SMALLINT,
    last_email TEXT,
    verified_email TEXT,
    last_phone TEXT,
    verified_phone TEXT,
    sequence BIGINT,
    password_set BOOLEAN,
    login_names TEXT,
    preferred_login_name TEXT,

    PRIMARY KEY (id)
);


CREATE TABLE adminapi.orgs (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    org_state SMALLINT,
    sequence BIGINT,

    domain TEXT,
    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE adminapi.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE adminapi.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE adminapi.current_sequences (
    view_name TEXT,
    timestamp TIMESTAMPTZ,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE adminapi.iam_members (
    user_id TEXT,

    iam_id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    user_name TEXT,
    email_address TEXT,
    first_name TEXT,
    last_name TEXT,
    roles TEXT ARRAY,
    display_name TEXT,
    sequence BIGINT,

    PRIMARY KEY (user_id)
);

CREATE TABLE management.orgs (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    org_state SMALLINT,
    sequence BIGINT,

    domain TEXT,
    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE management.org_members (
    user_id TEXT,
    org_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    user_name TEXT,
    email_address TEXT,
    first_name TEXT,
    last_name TEXT,
    roles TEXT ARRAY,
    display_name TEXT,
    sequence BIGINT,

    PRIMARY KEY (org_id, user_id)
);


CREATE TABLE auth.keys (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    resource_owner TEXT,
    private BOOLEAN,
    expiry TIMESTAMPTZ,
    algorithm TEXT,
    usage SMALLINT,
    key JSONB,
    sequence BIGINT,

    PRIMARY KEY (id, private)
);

CREATE TABLE auth.applications (
     id TEXT,

     creation_date TIMESTAMPTZ,
     change_date TIMESTAMPTZ,
     sequence BIGINT,

     app_state SMALLINT,
     resource_owner TEXT,
     app_name TEXT,
     project_id TEXT,
     app_type SMALLINT,
     is_oidc BOOLEAN,
     oidc_client_id TEXT,
     oidc_redirect_uris TEXT ARRAY,
     oidc_response_types SMALLINT ARRAY,
     oidc_grant_types SMALLINT ARRAY,
     oidc_application_type SMALLINT,
     oidc_auth_method_type SMALLINT,
     oidc_post_logout_redirect_uris TEXT ARRAY,

     PRIMARY KEY (id)
);

CREATE TABLE auth.user_grants (
    id TEXT,
    resource_owner TEXT,
    project_id TEXT,
    user_id TEXT,
    org_name TEXT,
    project_name TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    display_name TEXT,
    email TEXT,
    role_keys TEXT Array,
    grant_id TEXT,

    grant_state SMALLINT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

CREATE TABLE auth.orgs (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    org_state SMALLINT,
    sequence BIGINT,

    domain TEXT,
    name TEXT,

    PRIMARY KEY (id)
);

CREATE TABLE authz.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE authz.current_sequences (
    view_name TEXT,
    timestamp TIMESTAMPTZ,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE authz.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE authz.user_grants (
    id TEXT,
    resource_owner TEXT,
    project_id TEXT,
    user_id TEXT,
    org_name TEXT,
    project_name TEXT,
    user_name TEXT,
    first_name TEXT,
    last_name TEXT,
    display_name TEXT,
    email TEXT,
    role_keys TEXT Array,
    grant_id TEXT,

    grant_state SMALLINT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (id)
);

CREATE TABLE authz.applications (
    id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    app_state SMALLINT,
    resource_owner TEXT,
    app_name TEXT,
    project_id TEXT,
    app_type SMALLINT,
    is_oidc BOOLEAN,
    oidc_client_id TEXT,
    oidc_redirect_uris TEXT ARRAY,
    oidc_response_types SMALLINT ARRAY,
    oidc_grant_types SMALLINT ARRAY,
    oidc_application_type SMALLINT,
    oidc_auth_method_type SMALLINT,
    oidc_post_logout_redirect_uris TEXT ARRAY,

    PRIMARY KEY (id)
);

CREATE TABLE authz.orgs (
    id TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner TEXT,
    org_state SMALLINT,
    sequence BIGINT,

    domain TEXT,
    name TEXT,

    PRIMARY KEY (id)
);
`
