CREATE TABLE zitadel.projections.password_age_policies (
    id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    sequence INT8 NULL,
    state INT2 NULL,
    resource_owner TEXT,
    
    is_default BOOLEAN,
    max_age_days INT8 NULL,
    expire_warn_days INT8 NULL,

    PRIMARY KEY (id)
);