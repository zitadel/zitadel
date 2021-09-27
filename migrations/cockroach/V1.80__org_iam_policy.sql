CREATE TABLE zitadel.projections.org_iam_policies (
    id STRING NOT NULL, --TODO: pk
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    resource_owner TEXT,
    
    is_default BOOLEAN,
    user_login_must_be_domain BOOLEAN
);