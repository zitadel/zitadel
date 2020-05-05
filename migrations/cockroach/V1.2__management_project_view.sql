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

COMMIT;
