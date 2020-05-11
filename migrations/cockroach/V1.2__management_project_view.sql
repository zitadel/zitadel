BEGIN;


CREATE TABLE management.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ,
    object_type TEXT,

    PRIMARY KEY (object_type)
);

CREATE TABLE management.current_sequences (
    view_name TEXT,

    current_sequence BIGINT,

    PRIMARY KEY (view_name)
);

CREATE TABLE management.failed_event (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE management.granted_projects (
    project_id TEXT,
    org_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,

    project_name TEXT,
    org_name TEXT,
    org_domain TEXT,
    project_type SMALLINT,
    project_state SMALLINT,
    resource_owner TEXT,
    grant_id TEXT,
    granted_role_keys TEXT Array,
    sequence BIGINT,


    PRIMARY KEY (project_id, org_id)
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

COMMIT;
