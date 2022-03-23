CREATE TABLE zitadel.projections.keys_certificate (
                                                 id STRING REFERENCES zitadel.projections.keys ON DELETE NO ACTION,
                                                 expiry TIMESTAMPTZ NOT NULL,
                                                 key BYTES NOT NULL,
                                                 certificate BYTES NOT NULL,

                                                 PRIMARY KEY (id)
);

CREATE TABLE zitadel.projections.apps_saml_configs(
                                                      app_id STRING REFERENCES zitadel.projections.apps (id) ON DELETE CASCADE,

                                                      entity_id STRING NOT NULL,
                                                      metadata STRING,
                                                      metadata_url STRING,

                                                      PRIMARY KEY (app_id)
);