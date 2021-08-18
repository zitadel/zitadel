ALTER TABLE auth.idp_configs ADD COLUMN oauth_authorization_endpoint TEXT;
ALTER TABLE adminapi.idp_configs ADD COLUMN oauth_authorization_endpoint TEXT;
ALTER TABLE management.idp_configs ADD COLUMN oauth_authorization_endpoint TEXT;

ALTER TABLE auth.idp_configs ADD COLUMN oauth_token_endpoint TEXT;
ALTER TABLE adminapi.idp_configs ADD COLUMN oauth_token_endpoint TEXT;
ALTER TABLE management.idp_configs ADD COLUMN oauth_token_endpoint TEXT;