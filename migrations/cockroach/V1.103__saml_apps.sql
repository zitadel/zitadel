
CREATE TABLE zitadel.projections.apps_saml_configs(
                                                      app_id STRING REFERENCES zitadel.projections.apps (id) ON DELETE CASCADE,

                                                      entity_id STRING NOT NULL,
                                                      metadata STRING,
                                                      metadata_url STRING,

                                                      PRIMARY KEY (app_id)
);