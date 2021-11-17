CREATE TABLE test.projections.org_members (
    org_id STRING NOT NULL
    , user_id STRING NOT NULL
    , roles STRING[]

    , PRIMARY KEY (org_id, user_id)
);

CREATE TABLE test.projections.iam_members (
    iam_id STRING NOT NULL
    , user_id STRING NOT NULL
    , roles STRING[]

    , PRIMARY KEY (iam_id, user_id)
);

CREATE TABLE test.projections.project_members (
    project_id STRING NOT NULL
    , user_id STRING NOT NULL
    , roles STRING[]

    , PRIMARY KEY (project_id, user_id)
);
