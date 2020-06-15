BEGIN;

DROP TABLE management.granted_projects;

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
    org_domain TEXT,
    project_state SMALLINT,
    resource_owner TEXT,
    org_id TEXT,
    granted_role_keys TEXT Array,
    sequence BIGINT,

    PRIMARY KEY (grant_id)
);

COMMIT;
