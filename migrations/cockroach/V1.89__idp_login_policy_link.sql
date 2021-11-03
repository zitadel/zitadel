CREATE TABLE zitadel.projections.idp_login_policy_links(
    idp_id STRING,
    aggregate_id STRING,
    provider_type INT2,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence INT8,
    resource_owner STRING,

    PRIMARY KEY (aggregate_id, idp_id)
);
