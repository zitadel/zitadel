package scripts

const V115IdpProviders = `
ALTER TABLE adminapi.idp_providers ADD COLUMN idp_state SMALLINT;
ALTER TABLE management.idp_providers ADD COLUMN idp_state SMALLINT;
ALTER TABLE auth.idp_providers ADD COLUMN idp_state SMALLINT;

`
