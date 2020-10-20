ALTER TABLE management.idp_configs DROP COLUMN logo_src;
ALTER TABLE adminapi.idp_configs DROP COLUMN logo_src;
ALTER TABLE auth.idp_configs DROP COLUMN logo_src;

ALTER TABLE management.idp_configs ADD COLUMN styling_type SMALLINT;
ALTER TABLE adminapi.idp_configs ADD COLUMN styling_type SMALLINT;
ALTER TABLE auth.idp_configs ADD COLUMN styling_type SMALLINT;

ALTER TABLE management.idp_providers ADD COLUMN styling_type SMALLINT;
ALTER TABLE adminapi.idp_providers ADD COLUMN styling_type SMALLINT;
ALTER TABLE auth.idp_providers ADD COLUMN styling_type SMALLINT;

