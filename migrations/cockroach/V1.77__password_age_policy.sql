CREATE TABLE zitadel.projections.password_age_policies (
    id STRING NOT NULL, --TODO: pk
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    
    max_age_days INT8 NULL,
    expire_warn_days INT8 NULL
);