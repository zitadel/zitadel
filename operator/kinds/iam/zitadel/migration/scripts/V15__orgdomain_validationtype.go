package scripts

const V15OrgDomainValidationType = `
BEGIN;

ALTER TABLE management.org_domains ADD COLUMN validation_type SMALLINT;

CREATE TABLE adminapi.users (
  id TEXT,

  creation_date TIMESTAMPTZ,
  change_date TIMESTAMPTZ,

  resource_owner TEXT,
  user_state SMALLINT,
  last_login TIMESTAMPTZ,
  password_change TIMESTAMPTZ,
  user_name TEXT,
  login_names TEXT ARRAY,
  preferred_login_name TEXT,
  first_name TEXT,
  last_name TEXT,
  nick_Name TEXT,
  display_name TEXT,
  preferred_language TEXT,
  gender SMALLINT,
  email TEXT,
  is_email_verified BOOLEAN,
  phone TEXT,
  is_phone_verified BOOLEAN,
  country TEXT,
  locality TEXT,
  postal_code TEXT,
  region TEXT,
  street_address TEXT,
  otp_state SMALLINT,
  sequence BIGINT,
  password_set BOOLEAN,
  password_change_required BOOLEAN,
  mfa_max_set_up SMALLINT,
  mfa_init_skipped TIMESTAMPTZ,
  init_required BOOLEAN,

  PRIMARY KEY (id)
);

COMMIT;`
