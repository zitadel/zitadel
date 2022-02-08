CREATE TABLE zitadel.projections.iam (
    id STRING NOT NULL
    , change_date TIMESTAMPTZ NOT NULL
    , sequence INT8 NOT NULL

    , global_org_id STRING DEFAULT ''
    , iam_project_id STRING DEFAULT ''
    , setup_started SMALLINT DEFAULT 0
    , setup_done SMALLINT DEFAULT 0

    , PRIMARY KEY (id)
);
