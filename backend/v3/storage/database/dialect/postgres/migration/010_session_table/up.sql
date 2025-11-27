CREATE TABLE zitadel.session_user_agents (
    instance_id TEXT
    , fingerprint_id TEXT CHECK (fingerprint_id <> '')
    , ip INET
    , description TEXT
    , headers JSONB

    , PRIMARY KEY (instance_id, fingerprint_id)
);

CREATE TABLE zitadel.sessions (
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL CHECK (id <> '')
    , token TEXT
    , user_agent_id TEXT
    , lifetime INTERVAL
    , expiration TIMESTAMPTZ
    , user_id TEXT -- this column in used for referential integrity
    , creator_id TEXT
    , created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
    , updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
--     , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, user_agent_id) REFERENCES zitadel.session_user_agents(instance_id, fingerprint_id) ON DELETE SET NULL (user_agent_id)
);

CREATE OR REPLACE FUNCTION update_expiration()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.expiration := NEW.updated_at + NEW.lifetime;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_expiration
    BEFORE INSERT OR UPDATE OF updated_at, lifetime ON zitadel.sessions
    FOR EACH ROW
EXECUTE FUNCTION update_expiration();

-- CREATE INDEX idx_sessions_token ON zitadel.sessions (instance_id, token); --TODO: sha256

-- CREATE OR REPLACE FUNCTION zitadel.set_session_user()
--     RETURNS TRIGGER AS $$
-- BEGIN
--     -- this should mostly be a noop because the user_id should not change
--     UPDATE zitadel.sessions SET user_id = NEW.payload->>'userId' WHERE instance_id = NEW.instance_id AND id = NEW.session_id AND user_id <> NEW.payload->>'userId';
--
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

CREATE TYPE zitadel.session_factor_type AS ENUM (
    'user',
    'password',
    'passkey', -- is also a challenge
    'identity_provider_intent',
    'totp',
    'otp_sms', -- is also a challenge
    'otp_email' -- is also a challenge
);

CREATE TABLE zitadel.session_factors (
    instance_id TEXT NOT NULL
    , session_id TEXT NOT NULL
    , type zitadel.session_factor_type NOT NULL
    , last_challenged_at TIMESTAMPTZ -- this is only set if the type is a challenge
    , challenged_payload JSONB
    , last_verified_at TIMESTAMPTZ
    , verified_payload JSONB

    , PRIMARY KEY (instance_id, session_id, type)
    , FOREIGN KEY (instance_id, session_id) REFERENCES zitadel.sessions(instance_id, id) ON DELETE CASCADE
);

-- CREATE TRIGGER trg_sync_session_user
--     AFTER INSERT OR UPDATE OR DELETE ON zitadel.session_factors
--     FOR EACH ROW
--     WHEN (NEW.type = 'user' AND NEW.last_verified_at > OLD.last_verified_at)
-- EXECUTE FUNCTION zitadel.set_session_user();

CREATE TABLE zitadel.session_metadata (
    instance_id TEXT NOT NULL
    , session_id TEXT NOT NULL
    , key TEXT NOT NULL CHECK (key <> '')
    , value BYTEA NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    , PRIMARY KEY (instance_id, session_id, key)

    , CONSTRAINT fk_session_metadata_session FOREIGN KEY (instance_id, session_id) REFERENCES zitadel.sessions (instance_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_session_metadata_key ON zitadel.session_metadata (key);
CREATE INDEX idx_session_metadata_value ON zitadel.session_metadata (sha256(value));

-- TODO(adlerhurst): these indexes can currently not be used by Postgres, because of the type conversion
-- the value can be a json but doesn't have to be.
-- CREATE INDEX idx_session_metadata_value_number ON zitadel.session_metadata ((value::NUMERIC)) WHERE jsonb_typeof(value) = 'number';
-- CREATE INDEX idx_session_metadata_value_string ON zitadel.session_metadata ((value#>>'{}')) WHERE jsonb_typeof(value) = 'string';
-- CREATE INDEX idx_session_metadata_value_boolean ON zitadel.session_metadata ((value::BOOLEAN)) WHERE jsonb_typeof(value) = 'boolean';

CREATE TRIGGER trg_set_updated_at_session_metadata
    BEFORE INSERT OR UPDATE ON zitadel.session_metadata
    FOR EACH ROW
    WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();