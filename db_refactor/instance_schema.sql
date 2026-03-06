DROP SCHEMA IF EXISTS zitadel CASCADE;
CREATE SCHEMA zitadel;

-- Languages
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.languages CASCADE;
CREATE TABLE zitadel.languages (
  code VARCHAR(2) NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  
  deleted_at TIMESTAMP DEFAULT NULL
);
INSERT INTO zitadel.languages (code, name)
VALUES 
  ('bg','bulgerian'),
  ('cs','czechian'),
  ('de','german'),
  ('en','english'),
  ('es','spanish'),
  ('fr','french'),
  ('hu','hungarian'),
  ('id','indonesian'),
  ('it','italian'),
  ('ja','japanese'),
  ('ko','korean'),
  ('mk','macedonian'),
  ('nl','dutch'),
  ('pl','polish'),
  ('pt','portuguese'),
  ('ro','romanian'),
  ('ru','russian'),
  ('sv','swedish'),
  ('zh','zhuang');

--------------------------------------- instance ---------------------------------------
-- instances
DROP TABLE IF EXISTS zitadel.instances CASCADE;
CREATE TABLE zitadel.instances (
  id VARCHAR(100) NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  default_org_id VARCHAR(100) NOT NULL,
  iam_project_id VARCHAR(100) NOT NULL,
  console_client_id VARCHAR(100) NOT NULL,
  console_app_id VARCHAR(100) NOT NULL,
  language VARCHAR(2) REFERENCES languages(code),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Other -----------------------------------------------------------------------------------------------------------------------------------------------------------

-- Security Polcies
-----------------------------------------------------------------------------------------------------------------------------------------------------------
-- NOTE: this is 'Secret Settings' in the UI
DROP TABLE IF EXISTS zitadel.security_policies CASCADE;
CREATE TABLE zitadel.security_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  enable_iframe_embedding BOOLEAN DEFAULT FALSE NOT NULL,
  enable_impersonation BOOLEAN DEFAULT FALSE NOT NULL,

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL,
  origins text[]
);

-- Security Generator
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TYPE IF EXISTS zitadel.secret_generator_types CASCADE;
CREATE TYPE zitadel.secret_generator_types AS ENUM (
  'INIT_CODE',
  'VERIFY_EMAIL_CODE',
  'VERIFY_PHONE_CODE',
  'PASSWORD_RESET_CODE',
  'PASSWORDLESS_INIT_CODE',
  'APP_SECRET',
  'OTP_SMS',
  'OTP_EMAIL'
);

DROP TABLE IF EXISTS zitadel.secret_generators CASCADE;
CREATE TABLE zitadel.secret_generators (
  id SERIAL NOT NULL PRIMARY KEY,
  -- TODO investigate if this will only ever be applied to instance
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  generator_type zitadel.secret_generator_types NOT NULL,
  length BIGINT NOT NULL,
  expiry BIGINT NOT NULL,
  include_lower_letters BOOLEAN NOT NULL,
  include_upper_letters BOOLEAN NOT NULL,
  include_digits BOOLEAN NOT NULL,
  include_symbols BOOLEAN NOT NULL,

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);


-- Ocid Settings
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.oidc_settings CASCADE;
CREATE TABLE zitadel.oidc_settings (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  access_token_lifetime BIGINT NOT NULL,
  id_token_lifetime BIGINT NOT NULL,
  refresh_token_idle_expiration BIGINT NOT NULL,
  refresh_token_expiration BIGINT NOT NULL,

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Privacy Policy
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TYPE IF EXISTS zitadel.privacy_policy_state CASCADE;
CREATE TYPE zitadel.privacy_policy_state AS ENUM (
	'POLICYSTATEUNSPECIFIED',
	'POLICYSTATEACTIVE',
	'POLICYSTATEREMOVED',

	'POLICYSTATECOUNT'
);

-- NOTE this is 'External Links' in the UI
DROP TABLE IF EXISTS zitadel.privacy_policies CASCADE;
CREATE TABLE zitadel.privacy_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  -- https://zitadel.slack.com/archives/C07SVSZU38X/p1744291107892309
  state privacy_policy_state NOT NULL,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  privacy_link VARCHAR(200) NOT NULL,
  tos_link VARCHAR(200) NOT NULL,
  help_link VARCHAR(200) NOT NULL,
  support_email VARCHAR(200) NOT NULL,
  docs_link VARCHAR(200) DEFAULT 'https://zitadel.com/docs'::text NOT NULL,
  custom_link VARCHAR(200) NOT NULL,
  custom_link_text VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Apperance -----------------------------------------------------------------------------------------------------------------------------------------------------------

-- Custom Text 
-----------------------------------------------------------------------------------------------------------------------------------------------------------
-- 'Login Interface Text' in the UI
DROP TABLE IF EXISTS zitadel.custom_texts CASCADE;
CREATE TABLE zitadel.custom_texts (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,

  is_default BOOLEAN NOT NULL,
  template VARCHAR(200) NOT NULL,
  language VARCHAR(2) REFERENCES languages(code),
  key VARCHAR(30) NOT NULL,
  text VARCHAR(200) NOT NULL,


  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Message Text
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.message_text CASCADE;
CREATE TABLE zitadel.message_text (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  -- NOTE in the code the state field is the same as in privacy polciies
  state privacy_policy_state NOT NULL,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  type text NOT NULL,
  language VARCHAR(2) REFERENCES languages(code),
  title VARCHAR(100),
  pre_header VARCHAR(100),
  subject VARCHAR(100),
  greeting VARCHAR(100),
  text VARCHAR(100),
  button_text VARCHAR(100),
  footer_text VARCHAR(100),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Styling ; 'Branding' in the UI
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TYPE IF EXISTS zitadel.label_policy_state CASCADE;
CREATE TYPE zitadel.label_policy_state AS ENUM (
	'LabelPolicyStateUnspecified',
	'LabelPolicyStateActive',
	'LabelPolicyStateRemoved',
	'LabelPolicyStatePreview',

	'labelPolicyStateCount'
);
-- NOTE: styling is in adminapi schema currently NOT adminapi.styling2
DROP TABLE IF EXISTS zitadel.styling CASCADE;
CREATE TABLE zitadel.styling (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  -- label_policy_state label_policy_state NOT NULL,
  state label_policy_state NOT NULL,
  primary_color TEXT,
  background_color TEXT,
  warn_color TEXT,
  font_color TEXT,
  primary_color_dark TEXT,
  background_color_dark TEXT,
  warn_color_dark TEXT,
  font_color_dark TEXT,
  logo_url VARCHAR(200),
  icon_url VARCHAR(200),
  logo_dark_url VARCHAR(200),
  icon_dark_url VARCHAR(200),
  font_url VARCHAR(200),
  err_msg_popup BOOLEAN,
  disable_watermark BOOLEAN,
  hide_login_name_suffix BOOLEAN,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Domain Policies
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TYPE IF EXISTS zitadel.policy_state CASCADE;
CREATE TYPE zitadel.policy_state AS ENUM (
	'PolicyStateUnspecified',
	'PolicyStateActive',
	'PolicyStateRemoved',

	'policyStateCount'
);
DROP TABLE IF EXISTS zitadel.domain_policies CASCADE;
CREATE TABLE zitadel.domain_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  state policy_state NOT NULL,
  user_login_must_be_domain BOOLEAN NOT NULL,
  validate_org_domains BOOLEAN NOT NULL,
  smtp_sender_address_matches_instance_domain BOOLEAN NOT NULL,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Password/Security
-----------------------------------------------------------------------------------------------------------------------------------------------------------
-- TODO consider combining lockout_policies + password_complexity_policies + password_age_policies
-- Lockout Policies
DROP TABLE IF EXISTS zitadel.lockout_policies CASCADE;
CREATE TABLE zitadel.lockout_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  state policy_state NOT NULL,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  max_password_attempts BIGINT NOT NULL,
  max_otp_attempts BIGINT DEFAULT 0 NOT NULL,
  show_failure BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Password Age
DROP TABLE IF EXISTS zitadel.password_age_policies CASCADE;
CREATE TABLE zitadel.password_age_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  state policy_state NOT NULL,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  expire_warn_days BIGINT NOT NULL,
  max_age_days BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- Password Complexity
DROP TABLE IF EXISTS zitadel.password_complexity_policies CASCADE;
CREATE TABLE zitadel.password_complexity_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  state policy_state NOT NULL,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  min_length BIGINT NOT NULL,
  has_lowercase BOOLEAN NOT NULL,
  has_uppercase BOOLEAN NOT NULL,
  has_symbol BOOLEAN NOT NULL,
  has_number BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- idp
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TYPE IF EXISTS zitadel.idp_type CASCADE;
CREATE TYPE zitadel.idp_type AS ENUM ('oauth2', 'oidc', 'jwt', 'azure', 'github', 'gitlab', 'github_enterprise', 'gitlab_self_hosted', 'google', 'ldap', 'saml', 'apple');
DROP TABLE IF EXISTS zitadel.idp_providers CASCADE;
CREATE TABLE zitadel.idp_providers (
  id SERIAL NOT NULL PRIMARY KEY,
  idp_type idp_type NOT NULL
);

DROP TABLE IF EXISTS zitadel.idp_saml CASCADE;
CREATE TABLE zitadel.idp_saml(
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  with_signed_request BOOLEAN,
  name_id_format SMALLINT,
  transient_mapping_attribute_name VARCHAR(100),
  metadata bytea NOT NULL,
-- TODO maybe key should be another type
  key jsonb NOT NULL,
  certificate bytea NOT NULL,
  binding text
);

DROP TABLE IF EXISTS zitadel.idp_apple CASCADE;
CREATE TABLE zitadel.idp_apple (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  team_id VARCHAR(10) NOT NULL,
  key_id VARCHAR(10) NOT NULL,
  -- TODO maybe private_key should be another type
  private_key jsonb NOT NULL,
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_ldap CASCADE;
CREATE TABLE zitadel.idp_ldap (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  -- TODO check ALL sizes
  start_tls BOOLEAN NOT NULL,
  base_dn VARCHAR(100) NOT NULL,
  bind_dn VARCHAR(100)  NOT NULL,
  bind_password VARCHAR(100) NOT NULL,
  user_base VARCHAR(100) NOT NULL,
  timeout BIGINT NOT NULL,
  id_attribute VARCHAR(50),
  first_name_attribute VARCHAR(50),
  last_name_attribute VARCHAR(50),
  display_name_attribute VARCHAR(100),
  nick_name_attribute  VARCHAR(100),
  preferred_username_attribute VARCHAR(100),
  email_attribute VARCHAR(100),
  email_verified VARCHAR(100),
  phone_attribute VARCHAR(100),
  phone_verified_attribute VARCHAR(100),
  preferred_language_attribute VARCHAR(100),
  avatar_url_attribute VARCHAR(100),
  profile_attribute VARCHAR(100),
  user_object_classes text[] NOT NULL,
  user_filters text[] NOT NULL,
  servers text[] NOT NULL,
  root_ca bytea
);

DROP TABLE IF EXISTS zitadel.idp_google CASCADE;
CREATE TABLE zitadel.idp_google (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB  NOT NULL,
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_gitlab_self_hosted CASCADE;
CREATE TABLE zitadel.idp_gitlab_self_hosted (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  issuer VARCHAR(200) NOT NULL,
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB  NOT NULL,
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_gitlab CASCADE;
CREATE TABLE zitadel.idp_gitlab (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB NOT NULL,
  -- TODO check if there's a better type for scopes than text[]
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_github_enterprise CASCADE;
CREATE TABLE zitadel.idp_github_enterprise (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB  NOT NULL,
  authorization_endpoint VARCHAR(200) NOT NULL,
  token_endpoint VARCHAR(200) NOT NULL,
  user_endpoint VARCHAR(100) NOT NULL,
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_github CASCADE;
CREATE TABLE zitadel.idp_github(
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB NOT NULL,
  scopes text[]
);

DROP TYPE IF EXISTS zitadel.tenant CASCADE;
CREATE TYPE zitadel.tenant AS ENUM ('consumer', 'common', 'organization');
DROP TABLE IF EXISTS zitadel.idp_azure CASCADE;
CREATE TABLE zitadel.idp_azure (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB  NOT NULL,
  tenant tenant NOT NULL,
  is_email_verified BOOLEAN DEFAULT FALSE NOT NULL,
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_jwt CASCADE;
CREATE TABLE zitadel.idp_jwt (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  issuer VARCHAR(200) NOT NULL,
  jwt_endpoint VARCHAR(200) NOT NULL,
  keys_endpoint VARCHAR(200) NOT NULL,
  header_name VARCHAR(100)
);

DROP TABLE IF EXISTS zitadel.idp_oidc CASCADE;
CREATE TABLE zitadel.idp_oidc (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  issuer VARCHAR(200) NOT NULL,
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB  NOT NULL,
  id_token_mapping BOOLEAN DEFAULT FALSE NOT NULL,
  use_pkce BOOLEAN DEFAULT FALSE NOT NULL,
  scopes text[]
);

DROP TABLE IF EXISTS zitadel.idp_oauth2 CASCADE;
CREATE TABLE zitadel.idp_oauth2 (
  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id),
  client_id VARCHAR(100) NOT NULL,
  client_secret JSONB  NOT NULL,
  authorization_endpoint VARCHAR(200) NOT NULL,
  token_endpoint VARCHAR(200) NOT NULL,
  user_endpoint VARCHAR(100)  NOT NULL,
  id_attribute VARCHAR(100) NOT NULL,
  use_pkce BOOLEAN DEFAULT FALSE NOT NULL,
  scopes text[]
);

DROP TYPE IF EXISTS zitadel.idps_config_state CASCADE;
CREATE TYPE zitadel.idps_config_state AS ENUM (
	'IDPConfigStateUnspecified',
	'IDPConfigStateActive',
	'IDPConfigStateInactive',
	'IDPConfigStateRemoved',

	'idpConfigStateCount'
);

DROP TABLE IF EXISTS zitadel.idps CASCADE;
CREATE TABLE zitadel.idps ( -- based on idp_templates6
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  state idps_config_state NOT NULL,
  name VARCHAR(100) NOT NULL,
  owner_type SMALLINT NOT NULL,
  type SMALLINT NOT NULL,
  is_creation_allowed BOOLEAN DEFAULT FALSE NOT NULL,
  is_linking_allowed BOOLEAN DEFAULT FALSE NOT NULL,
  is_auto_creation BOOLEAN DEFAULT FALSE NOT NULL,
  is_auto_update BOOLEAN DEFAULT FALSE NOT NULL,
  auto_linking SMALLINT DEFAULT 0 NOT NULL,

  idp_provider_id INTEGER NOT NULL,
  CONSTRAINT idp_provider_id_fk FOREIGN KEY(idp_provider_id) REFERENCES idp_providers(id) ON DELETE CASCADE,

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- login policies
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.login_policies CASCADE;
CREATE TABLE zitadel.login_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  is_default BOOLEAN DEFAULT FALSE NOT NULL,
  allow_register BOOLEAN NOT NULL,
  allow_username_password BOOLEAN NOT NULL,
  allow_external_idps BOOLEAN NOT NULL,
  force_mfa BOOLEAN NOT NULL,
  force_mfa_local_only BOOLEAN DEFAULT FALSE NOT NULL,
  second_factors SMALLINT[],
  multi_factors SMALLINT[],
  passwordless_type SMALLINT NOT NULL,
  hide_password_reset BOOLEAN NOT NULL,
  ignore_unknown_usernames BOOLEAN NOT NULL,
  allow_domain_discovery BOOLEAN NOT NULL,
  disable_login_with_email BOOLEAN NOT NULL,
  disable_login_with_phone BOOLEAN NOT NULL,
  default_redirect_uri VARCHAR(200),
  password_check_lifetime BIGINT NOT NULL,
  external_login_check_lifetime BIGINT NOT NULL,
  mfa_init_skip_lifetime BIGINT NOT NULL,
  second_factor_check_lifetime BIGINT NOT NULL,
  multi_factor_check_lifetime BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- TODO not sure what projections.restrictions2 is for
-- CREATE TABLE zitadel.restrictions (
--     aggregate_id text NOT NULL,
--     creation_date timestamp with time zone NOT NULL,
--     change_date timestamp with time zone NOT NULL,
--     resource_owner text NOT NULL,
--     instance_id text NOT NULL,
--     sequence BIGINT NOT NULL,
--     disallow_public_org_registration BOOLEAN,
--     allowed_languages text[]
-- );

-- sms
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TYPE IF EXISTS zitadel.sms_type CASCADE;
CREATE TYPE zitadel.sms_type AS ENUM ('sms', 'http', 'twilio');

DROP TABLE IF EXISTS zitadel.sms_configs CASCADE;
CREATE TABLE zitadel.sms_configs (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  sms_type sms_type NOT NULL,
  sms_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

DROP TYPE IF EXISTS zitadel.sms_state CASCADE;
CREATE TYPE zitadel.sms_state AS ENUM (
	'SMSConfigStateUnspecified',
	'SMSConfigStateActive',
	'SMSConfigStateInactive',
	'SMSConfigStateRemoved'
);

DROP TABLE IF EXISTS zitadel.sms CASCADE;
CREATE TABLE zitadel.sms (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk     FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  -- I'm assuming the aggregate_id will always be the instance_id, so I'm remocing it
  -- aggregate_id text NOT NULL,
  state sms_state NOT NULL,
  -- I'm assuming the resource_owner will always be the instance_id, so I'm remocing it
  -- resource_owner text NOT NULL,
  description VARCHAR(200)  NOT NULL
);

DROP TABLE IF EXISTS zitadel.sms_http CASCADE;
CREATE TABLE zitadel.sms_http (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  endpoint VARCHAR(200) NOT NULL
);

DROP TABLE IF EXISTS zitadel.sms_twilio CASCADE;
CREATE TABLE zitadel.sms_twilio (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk     FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  -- I'm skeptical about all these lenghts, they're from the GRPC defintiions
  -- but they seem way too big, TODO double check this
  sid VARCHAR(200) NOT NULL,
  sender_number VARCHAR(200) NOT NULL,
  token VARCHAR(200) NOT NULL,
  verify_service_sid VARCHAR(200) NOT NULL
);

-- notifications
-----------------------------------------------------------------------------------------------------------------------------------------------------------
-- TODO sort out notificaitons
DROP TYPE IF EXISTS zitadel.notifications_types CASCADE;
CREATE TYPE zitadel.notifications_types AS ENUM (
  'PASSWORD_CHANGE'
);

DROP TYPE IF EXISTS zitadel.notification_method CASCADE;
CREATE TYPE zitadel.notification_method AS ENUM (
  'NOTIFICATION_TYPE_Unspecified',
  'NOTIFICATION_TYPE_Email',
  'NOTIFICATION_TYPE_SMS'
);

DROP TABLE IF EXISTS zitadel.notification_policies CASCADE;
CREATE TABLE zitadel.notification_policies (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk     FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  state policy_state NOT NULL,
  is_default BOOLEAN NOT NULL,
  -- NOTE: discuss if this should or shouldn't be an enum
  notification_type notifications_types NOT NULL,
  notification_method notification_method NOT NULL,
  text TEXT NOT NULL,
  send_notification BOOLEAN NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- smtp config
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.smtp_configs CASCADE;
CREATE TABLE zitadel.smtp_configs (
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk     FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  tls BOOLEAN NOT NULL,
  sender_address VARCHAR(200) NOT NULL,
  sender_name VARCHAR(200) NOT NULL,
  reply_to_address VARCHAR(200) NOT NULL,
  host VARCHAR(200) NOT NULL,
  -- need to sign off on username legnth
  username VARCHAR(100) NOT NULL,
  -- need to sign off on password legnth/can password be null?
  password VARCHAR(100) 
);

-- instance features
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.instance_features CASCADE;
CREATE TABLE zitadel.instance_features(
  -- adding id to as a means of keeping an audit on changes
  id SERIAL NOT NULL PRIMARY KEY,
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk     FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  -- maybe key should be a enum?
  key VARCHAR(100) NOT NULL,
  -- value seems to be always bool, maybe this might be worth changing to bool
  value JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- organizations
-----------------------------------------------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS zitadel.organizations CASCADE;
CREATE TABLE zitadel.organizations(
  id VARCHAR(100) NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  instance_id VARCHAR(100) NOT NULL,
  -- TODO turn org_state into an enum
  org_state SMALLINT NOT NULL,
  -- TODO ADD DOMAIN
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

--------------------------------------- organizations ---------------------------------------
DROP TABLE IF EXISTS zitadel.domains CASCADE;
CREATE TABLE zitadel.domains(
  id SERIAL NOT NULL PRIMARY KEY,
  organization_id VARCHAR(100) NOT NULL,
  instance_id VARCHAR(100) NOT NULL,
  domain VARCHAR(200) NOT NULL,
  is_verified BOOLEAN NOT NULL,
  is_primary BOOLEAN NOT NULL,
  -- TODO make validation_type enum
  validation_type SMALLINT NOT NULL,
  is_owner_removed BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP DEFAULT NULL
);

-- organization <----> domains
DROP TABLE IF EXISTS zitadel.organizations_domains_map CASCADE;
CREATE TABLE zitadel.organizations_domains_map (
  organization_id VARCHAR(100) NOT NULL,
  CONSTRAINT organization_id_fk FOREIGN KEY(organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
  domain_id INTEGER NOT NULL,
  CONSTRAINT domain_id_fk       FOREIGN KEY(domain_id) REFERENCES domains(id) ON DELETE CASCADE
);

-- instances <----> organizations
DROP TABLE IF EXISTS zitadel.instances_orgs_map CASCADE;
CREATE TABLE zitadel.instances_orgs_map (
  instance_id VARCHAR(100) NOT NULL,
  CONSTRAINT instance_id_fk     FOREIGN KEY(instance_id) REFERENCES instances(id) ON DELETE CASCADE,
  organization_id VARCHAR(100) NOT NULL,
  CONSTRAINT organization_id_fk FOREIGN KEY(organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);
