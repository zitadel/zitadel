SET enable_experimental_alter_column_type_general = true;

ALTER TABLE zitadel.projections.apps_oidc_configs ALTER version TYPE INT2 USING version::INT2;