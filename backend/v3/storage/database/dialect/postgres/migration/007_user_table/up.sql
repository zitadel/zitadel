-- ------------------------------------------------------------
-- table definitions
-- ------------------------------------------------------------
CREATE TABLE zitadel.verifications(
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL

    , value TEXT
    , code BYTES

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , expiry INTERVAL

    , failed_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_attempts >= 0)

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION zitadel.cleanup_verification(instance_id TEXT, verification_id TEXT)
    RETURNS VOID AS $$
BEGIN
    DELETE FROM zitadel.verifications v WHERE v.instance_id = instance_id AND v.id = verification_id;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE zitadel.user_state AS ENUM (
    'initial'
    , 'active'
    , 'inactive'
    , 'locked'
    , 'suspended'
);

-- user
CREATE TABLE zitadel.users(
    instance_id TEXT NOT NULL
    , organization_id TEXT NOT NULL
    , id TEXT NOT NULL CHECK (id <> '')

    , username TEXT NOT NULL CHECK (username <> '')
    , username_org_unique BOOLEAN NOT NULL -- this field MUST be filled if the username must be unique on organization level
    , state zitadel.user_state NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    
    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations(instance_id, id)
);

CREATE UNIQUE INDEX ON zitadel.users(instance_id, organization_id, username) WHERE username_org_unique IS TRUE; --TODO(adlerhurst): does that work if a username is already present on a user without org unique?
CREATE UNIQUE INDEX ON zitadel.users(instance_id, username) WHERE username_org_unique IS FALSE;

CREATE TABLE zitadel.user_metadata (
    instance_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    , key TEXT NOT NULL CHECK (key <> '')
    , value BYTEA NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    
    , PRIMARY KEY (instance_id, user_id, key)
    
    , CONSTRAINT fk_user_metadata_org FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users (instance_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_user_metadata_key ON zitadel.user_metadata (key);
CREATE INDEX idx_user_metadata_value ON zitadel.user_metadata (sha256(value));

CREATE TRIGGER trg_set_updated_at_user_metadata
  BEFORE INSERT OR UPDATE ON zitadel.user_metadata
  FOR EACH ROW
  WHEN (NEW.updated_at IS NULL)
  EXECUTE FUNCTION zitadel.set_updated_at();

-- machine user
CREATE TABLE zitadel.machine_users(
    name TEXT NOT NULL CHECK (name <> '')
    , description TEXT
    , secret BYTES
    , access_token_type SMALLINT

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations
) INHERITS (zitadel.users);

CREATE UNIQUE INDEX ON zitadel.machine_users(instance_id, organization_id, username) WHERE username_org_unique IS TRUE; --TODO(adlerhurst): does that work if a username is already present on a user without org unique?
CREATE UNIQUE INDEX ON zitadel.machine_users(instance_id, username) WHERE username_org_unique IS FALSE;

CREATE INDEX idx_machine_name ON zitadel.machine_users (instance_id, name);
CREATE INDEX idx_machine_user_username ON zitadel.machine_users (instance_id, username);
CREATE INDEX idx_machine_user_username_insensitive ON zitadel.machine_users (instance_id, lower(username));

-- human user
CREATE TABLE zitadel.human_users(
    first_name TEXT CHECK (first_name <> '')
    , last_name TEXT CHECK (last_name <> '')
    , nickname TEXT
    , display_name TEXT CHECK (display_name   <> '')
    , preferred_language TEXT CHECK (preferred_language <> '')
    , gender SMALLINT 
    , avatar_key TEXT
    , multi_factor_initialization_skipped_at TIMESTAMPTZ

    , password BYTES
    , password_change_required BOOLEAN NOT NULL DEFAULT FALSE
    , password_verified_at TIMESTAMPTZ
    , unverified_password_id TEXT
    , failed_password_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_password_attempts >= 0)

    , email TEXT
    , unverified_email_id TEXT
    , email_verified_at TIMESTAMPTZ
    , email_otp_verification_id TEXT

    , phone TEXT
    , unverified_phone_id TEXT
    , phone_verified_at TIMESTAMPTZ
    , phone_otp_verification_id TEXT

    , PRIMARY KEY (instance_id, organization_id, id)
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations
    , FOREIGN KEY (instance_id, unverified_password_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_password_id)
    , FOREIGN KEY (instance_id, unverified_email_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_email_id)
    , FOREIGN KEY (instance_id, unverified_phone_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_phone_id)
    , FOREIGN KEY (instance_id, email_otp_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (email_otp_verification_id)
    , FOREIGN KEY (instance_id, phone_otp_verification_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (phone_otp_verification_id)

    , UNIQUE (password_verification_id) WHERE password_verification_id IS NOT NULL
    , UNIQUE (email_verification_id) WHERE email_verification_id IS NOT NULL
    , UNIQUE (phone_verification_id) WHERE phone_verification_id IS NOT NULL
    , UNIQUE (email_otp_verification_id) WHERE email_otp_verification_id IS NOT NULL
    , UNIQUE (phone_otp_verification_id) WHERE phone_otp_verification_id IS NOT NULL

) INHERITS (zitadel.users);

CREATE UNIQUE INDEX ON zitadel.human_users(instance_id, organization_id, username) WHERE username_org_unique IS TRUE; --TODO(adlerhurst): does that work if a username is already present on a user without org unique?
CREATE UNIQUE INDEX ON zitadel.human_users(instance_id, username) WHERE username_org_unique IS FALSE;

CREATE INDEX idx_human_user_username ON zitadel.human_users (instance_id, username);
CREATE INDEX idx_human_user_username_insensitive ON zitadel.human_users (instance_id, lower(username));

CREATE OR REPLACE FUNCTION zitadel.cleanup_human_password_verification()
    RETURNS TRIGGER AS $$
BEGIN
    EXECUTE zitadel.cleanup_verification(OLD.instance_id, OLD.password_verification_id);
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cleanup_password_verification_update
    AFTER UPDATE ON users
    FOR EACH ROW
    WHEN (OLD.password_verification_id IS NOT NULL AND OLD.password_verification_id IS DISTINCT FROM NEW.password_verification_id)
EXECUTE FUNCTION cleanup_human_password_verification();

CREATE TRIGGER cleanup_password_verification_delete
    AFTER DELETE ON users
    FOR EACH ROW
    WHEN (OLD.password_verification_id IS NOT NULL)
EXECUTE FUNCTION cleanup_human_password_verification();

CREATE TABLE IF NOT EXISTS zitadel.user_personal_access_tokens(
    instance_id TEXT NOT NULL
    , token_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , user_id TEXT NOT NULL
    , expiration TIMESTAMPTZ
    , scopes TEXT[]
    
    , PRIMARY KEY (instance_id, token_id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

CREATE TYPE zitadel.passkey_type AS ENUM (
    'passwordless'
    , 'u2f'
);

CREATE TYPE zitadel.passkey_state AS ENUM (
    'unverified'
    , 'verified'
);

CREATE TABLE zitadel.human_passkeys(
    instance_id TEXT NOT NULL
    , token_id TEXT NOT NULL

    , user_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , type zitadel.passkey_type NOT NULL
    , name TEXT NOT NULL CHECK (name <> '')
    , sign_count INT NOT NULL DEFAULT 0 CHECK (sign_count >= 0)
    , challenge BYTES NOT NULL
    , is_verified BOOLEAN NOT NULL DEFAULT FALSE

    , PRIMARY KEY (instance_id, token_id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

CREATE TABLE zitadel.human_identity_provider_links(
    instance_id TEXT NOT NULL
    , identity_provider_id TEXT NOT NULL
    , provided_user_id TEXT NOT NULL
    , provided_username TEXT NOT NULL
    , user_id TEXT NOT NULL

    , PRIMARY KEY (instance_id, provided_user_id, identity_provider_id)

    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.human_users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, identity_provider_id) REFERENCES zitadel.identity_providers(instance_id, id) ON DELETE CASCADE
);