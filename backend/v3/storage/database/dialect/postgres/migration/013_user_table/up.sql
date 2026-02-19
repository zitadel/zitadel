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
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE

    -- human

    , first_name                            TEXT        CHECK ((type = 'machine' AND first_name IS NULL)                            OR (type = 'human' AND first_name <> ''))
    , last_name                             TEXT        CHECK ((type = 'machine' AND last_name IS NULL)                             OR (type = 'human' AND last_name <> ''))
    , nickname                              TEXT        CHECK ((type = 'machine' AND nickname IS NULL)                              OR (type = 'human'))
    , display_name                          TEXT        CHECK ((type = 'machine' AND display_name IS NULL)                          OR (type = 'human'))
    , preferred_language                    TEXT        CHECK ((type = 'machine' AND preferred_language IS NULL)                    OR (type = 'human'))
    , gender                                SMALLINT    CHECK ((type = 'machine' AND gender IS NULL)                                OR (type = 'human'))
    , avatar_key                            TEXT        CHECK ((type = 'machine' AND avatar_key IS NULL)                            OR (type = 'human'))
    , multifactor_initialization_skipped_at TIMESTAMPTZ CHECK ((type = 'machine' AND multifactor_initialization_skipped_at IS NULL) OR (type = 'human'))

    , password_hash                         TEXT        CHECK ((type = 'machine' AND password_hash IS NULL)                         OR (type = 'human'))
    , password_change_required              BOOLEAN     CHECK ((type = 'machine' AND password_change_required IS NULL)              OR (type = 'human'))
    , password_changed_at                   TIMESTAMPTZ CHECK ((type = 'machine' AND password_changed_at IS NULL)                   OR (type = 'human'))
    , password_verification_id              TEXT        CHECK ((type = 'machine' AND password_verification_id IS NULL)              OR (type = 'human')) -- used for reset password flow
    , password_last_successful_check        TIMESTAMPTZ CHECK ((type = 'machine' AND password_last_successful_check IS NULL)        OR (type = 'human'))
    , password_failed_attempts              SMALLINT    CHECK ((type = 'machine' AND password_failed_attempts IS NULL)              OR (type = 'human') AND (password_failed_attempts >= 0))

    , email                                 TEXT        CHECK ((type = 'machine' AND email IS NULL)                                 OR (type = 'human'))
    , unverified_email                      TEXT        CHECK ((type = 'machine' AND unverified_email IS NULL)                      OR (type = 'human' AND unverified_email <> '')) -- after successful verification this column is not cleared.
    , email_verified_at                     TIMESTAMPTZ CHECK ((type = 'machine' AND email_verified_at IS NULL)                     OR (type = 'human'))
    , email_verification_id                 TEXT        CHECK ((type = 'machine' AND email_verification_id IS NULL)                 OR (type = 'human' AND (email IS DISTINCT FROM unverified_email OR email_verification_id IS NULL)))
    , email_otp_enabled_at                  TIMESTAMPTZ CHECK ((type = 'machine' AND email_otp_enabled_at IS NULL)                  OR (type = 'human'))
    , email_otp_last_successful_check       TIMESTAMPTZ CHECK ((type = 'machine' AND email_otp_last_successful_check IS NULL)       OR (type = 'human'))
    , email_otp_failed_attempts             SMALLINT    CHECK ((type = 'machine' AND email_otp_failed_attempts IS NULL)             OR (type = 'human'))

    , phone                                 TEXT        CHECK ((type = 'machine' AND phone IS NULL)                                 OR (type = 'human'))
    , unverified_phone                      TEXT        CHECK ((type = 'machine' AND unverified_phone IS NULL)                      OR (type = 'human')) -- after successful verification this column is not cleared.
    , phone_verified_at                     TIMESTAMPTZ CHECK ((type = 'machine' AND phone_verified_at IS NULL)                     OR (type = 'human'))
    , phone_verification_id                 TEXT        CHECK ((type = 'machine' AND phone_verification_id IS NULL)                 OR (type = 'human' AND (phone IS DISTINCT FROM unverified_phone OR phone_verification_id IS NULL)))
    , sms_otp_enabled_at                    TIMESTAMPTZ CHECK ((type = 'machine' AND sms_otp_enabled_at IS NULL)                    OR (type = 'human'))
    , sms_otp_last_successful_check         TIMESTAMPTZ CHECK ((type = 'machine' AND sms_otp_last_successful_check IS NULL)         OR (type = 'human'))
    , sms_otp_failed_attempts               SMALLINT    CHECK ((type = 'machine' AND sms_otp_failed_attempts IS NULL)               OR (type = 'human'))

    , totp_secret                           BYTEA       CHECK ((type = 'machine' AND totp_secret IS NULL)                           OR (type = 'human'))
    , totp_verified_at                      TIMESTAMPTZ CHECK ((type = 'machine' AND totp_verified_at IS NULL)                      OR (type = 'human'))
    , totp_last_successful_check            TIMESTAMPTZ CHECK ((type = 'machine' AND totp_last_successful_check IS NULL)            OR (type = 'human'))
    , totp_failed_attempts                  SMALLINT    CHECK ((type = 'machine' AND totp_failed_attempts IS NULL)                  OR (type = 'human'))

    , invite_verification_id                TEXT        CHECK ((type = 'machine' AND invite_verification_id IS NULL)                OR (type = 'human'))
    , invite_accepted_at                    TIMESTAMPTZ CHECK ((type = 'machine' AND invite_accepted_at IS NULL)                    OR (type = 'human'))
    , invite_failed_attempts                SMALLINT    CHECK ((type = 'machine' AND invite_failed_attempts IS NULL)                OR (type = 'human'))

    -- foreign keys for verifications are created in the verification migration

    -- machine
    
    , name                                  TEXT        CHECK ((type = 'human' AND name IS NULL)                                  OR (type = 'machine'))
    , description                           TEXT        CHECK ((type = 'human' AND description IS NULL)                           OR (type = 'machine'))
    , secret                                TEXT        CHECK ((type = 'human' AND secret IS NULL)                                OR (type = 'machine'))
    , access_token_type                     SMALLINT    CHECK ((type = 'human' AND access_token_type IS NULL)                     OR (type = 'machine'))
);

-- previously created tables
ALTER TABLE zitadel.identity_provider_intents ADD CONSTRAINT fk_idp_intent_user FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users (instance_id, id) ON DELETE CASCADE;
ALTER TABLE zitadel.sessions ADD CONSTRAINT fk_session_user FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE;
ALTER TABLE zitadel.authorizations ADD CONSTRAINT fk_authorization_user FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users (instance_id, id) ON DELETE CASCADE;

-- user
CREATE UNIQUE INDEX ON zitadel.users(instance_id, organization_id, username) WHERE username_org_unique IS TRUE; --TODO(adlerhurst): does that work if a username is already present on a user without org unique?
CREATE UNIQUE INDEX ON zitadel.users(instance_id, username) WHERE username_org_unique IS FALSE;
CREATE INDEX idx_user_username ON zitadel.users (username);
CREATE INDEX idx_user_username_lower ON zitadel.users (lower(username));
CREATE INDEX idx_machine_name ON zitadel.users (name);
CREATE INDEX idx_human_email ON zitadel.users (email);
CREATE INDEX idx_human_email_lower ON zitadel.users (lower(email));
CREATE INDEX idx_human_phone ON zitadel.users (phone);
CREATE INDEX idx_human_phone_lower ON zitadel.users (lower(phone));

-- human
CREATE UNIQUE INDEX ON zitadel.users(password_verification_id) WHERE password_verification_id IS NOT NULL;
CREATE UNIQUE INDEX ON zitadel.users(email_verification_id) WHERE email_verification_id IS NOT NULL;
CREATE UNIQUE INDEX ON zitadel.users(phone_verification_id) WHERE phone_verification_id IS NOT NULL;

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
    , scopes TEXT[]
    
    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
);

-- ----------------------------------------------------------------
-- machine keys
-- ----------------------------------------------------------------

CREATE TABLE zitadel.machine_keys(
    instance_id TEXT NOT NULL
    , user_id TEXT NOT NULL

    , id TEXT NOT NULL
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

CREATE TABLE zitadel.user_passkeys(
    instance_id TEXT NOT NULL
    , token_id TEXT NOT NULL
    , key_id BYTEA NOT NULL

    , user_id TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , verified_at TIMESTAMPTZ

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

CREATE INDEX idx_user_passkeys_challenge ON zitadel.user_passkeys (sha256(challenge));
CREATE INDEX idx_user_passkeys_key_id ON zitadel.user_passkeys (key_id);
CREATE INDEX idx_user_passkeys_type ON zitadel.user_passkeys (type);

-- ----------------------------------------------------------------
-- identity provider links
-- ----------------------------------------------------------------

CREATE TABLE zitadel.user_identity_provider_links(
    instance_id TEXT NOT NULL
    , identity_provider_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    
    , provided_user_id TEXT NOT NULL CHECK(provided_user_id <> '')  
    , provided_username TEXT NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , PRIMARY KEY (instance_id, identity_provider_id, provided_user_id)

    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, identity_provider_id) REFERENCES zitadel.identity_providers(instance_id, id) ON DELETE CASCADE
);