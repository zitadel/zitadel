CREATE TYPE zitadel.user_state AS ENUM (
    'initial'
    , 'active'
    , 'inactive'
    , 'locked'
    , 'suspended'
);

CREATE TYPE zitadel.user_type AS ENUM (
    'human'
    , 'machine'
);

CREATE TABLE zitadel.users(
    instance_id TEXT NOT NULL
    , organization_id TEXT NOT NULL
    , id TEXT NOT NULL CHECK (id <> '')

    , username TEXT NOT NULL CHECK (username <> '')
    , username_org_unique BOOLEAN DEFAULT FALSE NOT NULL -- this field MUST be filled if the username must be unique on organization level
    , state zitadel.user_state NOT NULL DEFAULT 'active'
    , type zitadel.user_type NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations(instance_id, id)

    -- human

    , first_name TEXT
    , last_name TEXT
    , nickname TEXT
    , display_name TEXT
    , preferred_language TEXT
    , gender SMALLINT
    , avatar_key TEXT
    , multifactor_initialization_skipped_at TIMESTAMPTZ

    , password BYTEA
    , password_change_required BOOLEAN
    , password_verified_at TIMESTAMPTZ
    , unverified_password_id TEXT
    , failed_password_attempts SMALLINT

    , email TEXT
    , email_verified_at TIMESTAMPTZ
    , unverified_email_id TEXT
    , email_otp_enabled_at TIMESTAMPTZ
    , last_successful_email_otp_check TIMESTAMPTZ
    , email_otp_verification_id TEXT

    , phone TEXT
    , phone_verified_at TIMESTAMPTZ
    , unverified_phone_id TEXT
    , sms_otp_enabled_at TIMESTAMPTZ
    , last_successful_sms_otp_check TIMESTAMPTZ
    , sms_otp_verification_id TEXT

    , totp_secret_id TEXT -- reference to a verification that holds the secret
    , totp_verified_at TIMESTAMPTZ
    , unverified_totp_id TEXT -- reference to a verification that holds the new secret during change
    , last_successful_totp_check TIMESTAMPTZ

    -- foreign keys for verifications are created in the verification migration

    -- machine
    
    , name TEXT
    , description TEXT
    , secret BYTEA
    , access_token_type SMALLINT
);

-- user
CREATE UNIQUE INDEX ON zitadel.users(instance_id, organization_id, username) WHERE username_org_unique IS TRUE; --TODO(adlerhurst): does that work if a username is already present on a user without org unique?
CREATE UNIQUE INDEX ON zitadel.users(instance_id, username) WHERE username_org_unique IS FALSE;
CREATE INDEX idx_user_username ON zitadel.users (username);
CREATE INDEX idx_user_username_insensitive ON zitadel.users (lower(username));
CREATE INDEX idx_machine_name ON zitadel.users (name);
CREATE INDEX idx_human_email ON zitadel.users (email);
CREATE INDEX idx_human_email_lower ON zitadel.users (lower(email));
CREATE INDEX idx_human_phone ON zitadel.users (phone);
CREATE INDEX idx_human_phone_lower ON zitadel.users (lower(phone));

-- human

CREATE UNIQUE INDEX ON zitadel.users(unverified_password_id) WHERE unverified_password_id IS NOT NULL;
CREATE UNIQUE INDEX ON zitadel.users(unverified_email_id) WHERE unverified_email_id IS NOT NULL;
CREATE UNIQUE INDEX ON zitadel.users(unverified_phone_id) WHERE unverified_phone_id IS NOT NULL;
CREATE UNIQUE INDEX ON zitadel.users(email_otp_verification_id) WHERE email_otp_verification_id IS NOT NULL;
CREATE UNIQUE INDEX ON zitadel.users(sms_otp_verification_id) WHERE sms_otp_verification_id IS NOT NULL;

CREATE OR REPLACE FUNCTION zitadel.validate_human_user()
RETURNS TRIGGER AS $$
BEGIN
    -- Validate that machine-specific fields are NULL
    IF NEW.name IS NOT NULL OR NEW.description IS NOT NULL OR NEW.secret IS NOT NULL OR NEW.access_token_type IS NOT NULL THEN
        RAISE EXCEPTION 'Machine-specific fields must be NULL for human users';
    END IF;

    -- Validate non-empty string fields
    IF NEW.first_name IS NOT NULL AND NEW.first_name = '' THEN
        RAISE EXCEPTION 'first_name cannot be empty string';
    END IF;
    IF NEW.last_name IS NOT NULL AND NEW.last_name = '' THEN
        RAISE EXCEPTION 'last_name cannot be empty string';
    END IF;
    IF NEW.display_name IS NOT NULL AND NEW.display_name = '' THEN
        RAISE EXCEPTION 'display_name cannot be empty string';
    END IF;
    IF NEW.preferred_language IS NOT NULL AND NEW.preferred_language = '' THEN
        RAISE EXCEPTION 'preferred_language cannot be empty string';
    END IF;

    -- Validate numeric constraints
    IF NEW.failed_password_attempts IS NOT NULL AND NEW.failed_password_attempts < 0 THEN
        RAISE EXCEPTION 'failed_password_attempts must be >= 0';
    END IF;

    -- Set defaults
    IF NEW.password_change_required IS NULL THEN
        NEW.password_change_required := FALSE;
    END IF;
    IF NEW.failed_password_attempts IS NULL THEN
        NEW.failed_password_attempts := 0;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_human_user
BEFORE INSERT OR UPDATE ON zitadel.users
FOR EACH ROW
WHEN (NEW.type = 'human')
EXECUTE FUNCTION zitadel.validate_human_user();

CREATE FUNCTION zitadel.validate_machine_user()
RETURNS TRIGGER AS $$
BEGIN
    -- Validate that human-specific fields are NULL
    IF NEW.first_name IS NOT NULL OR NEW.last_name IS NOT NULL OR NEW.nickname IS NOT NULL 
       OR NEW.display_name IS NOT NULL OR NEW.preferred_language IS NOT NULL OR NEW.gender IS NOT NULL
       OR NEW.avatar_key IS NOT NULL OR NEW.multifactor_initialization_skipped_at IS NOT NULL
       OR NEW.password IS NOT NULL OR NEW.password_change_required IS NOT NULL 
       OR NEW.password_verified_at IS NOT NULL OR NEW.unverified_password_id IS NOT NULL
       OR NEW.failed_password_attempts IS NOT NULL OR NEW.email IS NOT NULL 
       OR NEW.email_verified_at IS NOT NULL OR NEW.unverified_email_id IS NOT NULL
       OR NEW.email_otp_enabled_at IS NOT NULL OR NEW.last_successful_email_otp_check IS NOT NULL
       OR NEW.email_otp_verification_id IS NOT NULL OR NEW.phone IS NOT NULL
       OR NEW.phone_verified_at IS NOT NULL OR NEW.unverified_phone_id IS NOT NULL
       OR NEW.sms_otp_enabled_at IS NOT NULL OR NEW.last_successful_sms_otp_check IS NOT NULL
       OR NEW.sms_otp_verification_id IS NOT NULL OR NEW.totp_secret_id IS NOT NULL
       OR NEW.totp_verified_at IS NOT NULL OR NEW.unverified_totp_id IS NOT NULL
       OR NEW.last_successful_totp_check IS NOT NULL THEN
        RAISE EXCEPTION 'Human-specific fields must be NULL for machine users';
    END IF;

    -- Validate name is not empty
    IF NEW.name IS NULL OR NEW.name = '' THEN
        RAISE EXCEPTION 'name cannot be empty string';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_machine_user
BEFORE INSERT OR UPDATE ON zitadel.users
FOR EACH ROW
WHEN (NEW.type = 'machine')
EXECUTE FUNCTION zitadel.validate_machine_user();

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.users
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

-- ----------------------------------------------------------------
-- user metadata
-- ----------------------------------------------------------------

CREATE TABLE zitadel.user_metadata (
    instance_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    , key TEXT NOT NULL CHECK (key <> '')
    , value BYTEA NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    
    , PRIMARY KEY (instance_id, user_id, key)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users (instance_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_user_metadata_key ON zitadel.user_metadata (key);
CREATE INDEX idx_user_metadata_value ON zitadel.user_metadata (sha256(value));

CREATE TRIGGER trg_set_updated_at_user_metadata
  BEFORE INSERT OR UPDATE ON zitadel.user_metadata
  FOR EACH ROW
  WHEN (NEW.updated_at IS NULL)
  EXECUTE FUNCTION zitadel.set_updated_at();

-- ----------------------------------------------------------------
-- personal access tokens
-- ----------------------------------------------------------------

CREATE TABLE zitadel.user_personal_access_tokens(
    instance_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    
    , id TEXT NOT NULL
    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expires_at TIMESTAMPTZ
    , type SMALLINT NOT NULL CHECK (type >= 0) --TODO(adlerhurst): remove column
    , public_key BYTEA NOT NULL --TODO(adlerhurst): remove column
    , scopes TEXT[]
    
    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

-- ----------------------------------------------------------------
-- machine keys
-- ----------------------------------------------------------------

CREATE TABLE zitadel.machine_keys(
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL

    , user_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expires_at TIMESTAMPTZ

    , type SMALLINT NOT NULL CHECK (type >= 0)
    , public_key BYTEA NOT NULL

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

-- ----------------------------------------------------------------
-- passkeys
-- ----------------------------------------------------------------

CREATE TYPE zitadel.passkey_type AS ENUM (
    'passwordless'
    , 'u2f'
);

CREATE TABLE zitadel.human_passkeys(
    instance_id TEXT NOT NULL
    , token_id TEXT NOT NULL
    , key_id TEXT NOT NULL

    , user_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , verified_at TIMESTAMPTZ
    , init_verification_id TEXT

    , type zitadel.passkey_type NOT NULL
    , name TEXT NOT NULL CHECK (name <> '')
    , sign_count INT NOT NULL DEFAULT 0 CHECK (sign_count >= 0)
    , challenge BYTEA NOT NULL
    , public_key BYTEA NOT NULL
    , attestation_type TEXT NOT NULL CHECK (attestation_type <> '')
    , authenticator_attestation_guid BYTEA NOT NULL
    , relying_party_id TEXT NOT NULL CHECK (relying_party_id <> '')

    , PRIMARY KEY (instance_id, token_id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

-- ----------------------------------------------------------------
-- identity provider links
-- ----------------------------------------------------------------

CREATE TABLE zitadel.human_identity_provider_links(
    instance_id TEXT NOT NULL
    , identity_provider_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    
    , provided_user_id TEXT NOT NULL
    , provided_username TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , PRIMARY KEY (instance_id, identity_provider_id, provided_user_id)

    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, identity_provider_id) REFERENCES zitadel.identity_providers(instance_id, id) ON DELETE CASCADE
);