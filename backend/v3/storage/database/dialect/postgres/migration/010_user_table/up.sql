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

    , first_name TEXT CHECK (first_name <> '' AND type = 'human')
    , last_name TEXT CHECK (last_name <> '' AND type = 'human')
    , nickname TEXT CHECK (type = 'human')
    , display_name TEXT CHECK (display_name   <> '' AND type = 'human')
    , preferred_language TEXT CHECK (preferred_language <> '' AND type = 'human')
    , gender SMALLINT CHECK (type = 'human')
    , avatar_key TEXT CHECK (type = 'human')
    , multi_factor_initialization_skipped_at TIMESTAMPTZ CHECK (type = 'human')

    , password BYTEA CHECK (type = 'human')
    , password_change_required BOOLEAN CHECK (type = 'human')
    , password_verified_at TIMESTAMPTZ CHECK (type = 'human')
    , unverified_password_id TEXT CHECK (type = 'human')
    , failed_password_attempts SMALLINT DEFAULT 0 CHECK (failed_password_attempts >= 0 AND type = 'human')

    , email TEXT CHECK (type = 'human')
    , email_verified_at TIMESTAMPTZ CHECK (type = 'human')
    , unverified_email_id TEXT CHECK (type = 'human')
    , email_otp_enabled_at TIMESTAMPTZ CHECK (type = 'human')
    , last_successful_email_otp_check TIMESTAMPTZ CHECK (type = 'human')
    , email_otp_verification_id TEXT CHECK (type = 'human')

    , phone TEXT CHECK (type = 'human')
    , phone_verified_at TIMESTAMPTZ CHECK (type = 'human')
    , unverified_phone_id TEXT CHECK (type = 'human')
    , sms_otp_enabled_at TIMESTAMPTZ CHECK (type = 'human')
    , last_successful_sms_otp_check TIMESTAMPTZ CHECK (type = 'human')
    , sms_otp_verification_id TEXT CHECK (type = 'human')

    , totp_secret_id TEXT CHECK (type = 'human') -- reference to a verification that holds the secret
    , totp_verified_at TIMESTAMPTZ CHECK (type = 'human')
    , unverified_totp_id TEXT CHECK (type = 'human') -- reference to a verification that holds the new secret during change
    , last_successful_totp_check TIMESTAMPTZ CHECK (type = 'human')

    , FOREIGN KEY (instance_id, unverified_password_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_password_id)
    , FOREIGN KEY (instance_id, unverified_email_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_email_id)
    , FOREIGN KEY (instance_id, unverified_phone_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_phone_id)
    , FOREIGN KEY (instance_id, email_otp_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (email_otp_verification_id)
    , FOREIGN KEY (instance_id, sms_otp_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (sms_otp_verification_id)

    -- machine
    
    , name TEXT CHECK (name <> '' AND type = 'machine')
    , description TEXT CHECK (type = 'machine')
    , secret BYTEA CHECK (type = 'machine')
    , access_token_type SMALLINT CHECK (type = 'machine')
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

CREATE FUNCTION zitadel.validate_human_user()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.password_change_required IS NULL THEN
        NEW.password_change_required := FALSE;
    END IF;
    IF NEW.failed_password_attempts IS NULL THEN
        NEW.failed_password_attempts := 0;
    END IF;
    IF NEW.email_otp_enabled IS NULL THEN
        NEW.email_otp_enabled := FALSE;
    END IF;
    IF NEW.sms_otp_enabled IS NULL THEN
        NEW.sms_otp_enabled := FALSE;
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
    IF NEW.name IS NULL THEN
        NEW.name := '';
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

CREATE OR REPLACE FUNCTION zitadel.cleanup_orphaned_user_verifications()
RETURNS TRIGGER AS $$
BEGIN
    -- This function cleans up verifications if the corresponding ID in the users table
    -- is changed, removed or if the user is deleted.

    IF OLD.unverified_password_id IS NOT NULL AND OLD.unverified_password_id IS DISTINCT FROM NEW.unverified_password_id THEN
        DELETE FROM zitadel.verifications WHERE instance_id = OLD.instance_id AND id = OLD.unverified_password_id;
    END IF;

    IF OLD.unverified_email_id IS NOT NULL AND OLD.unverified_email_id IS DISTINCT FROM NEW.unverified_email_id THEN
        DELETE FROM zitadel.verifications WHERE instance_id = OLD.instance_id AND id = OLD.unverified_email_id;
    END IF;

    IF OLD.unverified_phone_id IS NOT NULL AND OLD.unverified_phone_id IS DISTINCT FROM NEW.unverified_phone_id THEN
        DELETE FROM zitadel.verifications WHERE instance_id = OLD.instance_id AND id = OLD.unverified_phone_id;
    END IF;

    IF OLD.email_otp_verification_id IS NOT NULL AND OLD.email_otp_verification_id IS DISTINCT FROM NEW.email_otp_verification_id THEN
        DELETE FROM zitadel.verifications WHERE instance_id = OLD.instance_id AND id = OLD.email_otp_verification_id;
    END IF;

    IF OLD.phone_otp_verification_id IS NOT NULL AND OLD.phone_otp_verification_id IS DISTINCT FROM NEW.phone_otp_verification_id THEN
        DELETE FROM zitadel.verifications WHERE instance_id = OLD.instance_id AND id = OLD.phone_otp_verification_id;
    END IF;

    RETURN NEW; -- Return value is ignored for AFTER triggers.
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_cleanup_verifications
AFTER UPDATE OR DELETE ON zitadel.users
FOR EACH ROW
EXECUTE FUNCTION zitadel.cleanup_orphaned_user_verifications();

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
    , token_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , user_id TEXT NOT NULL
    , expiration TIMESTAMPTZ
    , scopes TEXT[]
    
    , PRIMARY KEY (instance_id, token_id)
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
    , FOREIGN KEY (instance_id, init_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (init_verification_id)
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

    , UNIQUE (instance_id, user_id, provided_user_id)
);