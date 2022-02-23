alter table auth.keys add column certificate boolean default false;
alter table auth.keys alter column certificate set not null;
alter table auth.keys alter primary key using columns( id, private, certificate);
drop index if exists auth.keys_id_private_key cascade;

CREATE TABLE zitadel.projections.apps_saml_configs(
                                                      app_id STRING REFERENCES zitadel.projections.apps (id) ON DELETE CASCADE,

                                                      entity_id STRING NOT NULL,
                                                      metadata STRING,
                                                      metadata_url STRING,

                                                      PRIMARY KEY (app_id)
);