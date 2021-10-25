CREATE TABLE zitadel.projections.apps(
    id STRING,
    project_id STRING NOT NULL,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    resource_owner STRING NOT NULL,
    state INT2,
    sequence INT8,
    
    name STRING,

    PRIMARY KEY (id),
    INDEX idx_project_id (project_id)
);

CREATE TABLE zitadel.projections.apps_api_configs(
	app_id STRING REFERENCES zitadel.projections.apps (id) ON DELETE CASCADE,
	
    client_id STRING NOT NULL,
	client_secret JSONB,
	auth_method INT2,

    PRIMARY KEY (app_id)
);

CREATE TABLE zitadel.projections.apps_oidc_configs(
    app_id STRING REFERENCES zitadel.projections.apps (id) ON DELETE CASCADE,
    
    version STRING NOT NULL,
    client_id STRING NOT NULL,
    client_secret JSONB,
    redirect_uris STRING[],
    response_types INT2[],
    grant_types INT2[],
    application_type INT2,
    auth_method_type INT2,
    post_logout_redirect_uris STRING[],
    is_dev_mode BOOLEAN,
    access_token_type INT2,
    access_token_role_assertion BOOLEAN,
    id_token_role_assertion BOOLEAN,
    id_token_userinfo_assertion BOOLEAN,
    clock_skew INT8,
    additional_origins STRING[],

    PRIMARY KEY (app_id)
);
