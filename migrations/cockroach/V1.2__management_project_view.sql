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

COMMIT;
