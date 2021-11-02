CREATE TABLE zitadel.projections.idp_user_links(
    idp_id STRING,
    user_id STRING,
    external_user_id STRING,
    display_name STRING,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence INT8,
    resource_owner STRING,

    PRIMARY KEY (idp_id, external_user_id),
    INDEX idx_user (user_id)
);
