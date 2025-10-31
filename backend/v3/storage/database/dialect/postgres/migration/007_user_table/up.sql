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

    , password BYTES
    , password_change_required BOOLEAN NOT NULL DEFAULT FALSE
    , password_verified_at TIMESTAMPTZ
    , unverified_password_id TEXT
    , failed_password_attempts SMALLINT NOT NULL DEFAULT 0 CHECK (failed_password_attempts >= 0)

    , email TEXT
    , unverified_email_id TEXT
    , email_verified_at TIMESTAMPTZ

    , phone TEXT
    , unverified_phone_id TEXT
    , phone_verified_at TIMESTAMPTZ

    , PRIMARY KEY (instance_id, organization_id, id)
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations
    , FOREIGN KEY (instance_id, unverified_password_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_password_id)
    , FOREIGN KEY (instance_id, unverified_email_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_email_id)
    , FOREIGN KEY (instance_id, unverified_phone_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL (unverified_phone_id)

    , UNIQUE (password_verification_id) WHERE password_verification_id IS NOT NULL
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






































CREATE TYPE zitadel.human_contact_type AS ENUM (
    'email'
    , 'phone'
);

CREATE INDEX idx_human_contacts_value ON zitadel.human_contacts(instance_id, value);
CREATE INDEX idx_human_contacts_value_lower ON zitadel.human_contacts(instance_id, lower(value)) WHERE type = 'email';

CREATE TABLE zitadel.human_contacts(
    instance_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    
    , type zitadel.human_contact_type NOT NULL
    , value TEXT
    , verified_at TIMESTAMPTZ
    , unverified_value_id TEXT

    , PRIMARY KEY (instance_id, user_id, type)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.human_users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, unverified_value_id) REFERENCES zitadel.verifications(instance_id, id) ON DELETE SET NULL
);


-- CREATE TABLE zitadel.human_security(
--     instance_id TEXT NOT NULL
--     , user_id TEXT NOT NULL

--     , password_change_required BOOLEAN NOT NULL DEFAULT FALSE
--     , password_changed TIMESTAMPTZ
--     , mfa_init_skipped BOOLEAN NOT NULL DEFAULT FALSE

--     , PRIMARY KEY (instance_id, user_id)
--     , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.human_users(instance_id, id) ON DELETE CASCADE
-- );

-- ------------------------------------------------------------
-- function definitions
-- ------------------------------------------------------------

-- sets the username uniqueness initially
-- CREATE OR REPLACE FUNCTION zitadel.user_set_username_uniqueness()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     SELECT 
--         payload->'organizationScopedUsernames'::BOOLEAN INTO NEW.username_org_unique 
--     FROM 
--         zitadel.settings 
--     WHERE 
--         ((instance_id = NEW.instance_id AND organization_id = NEW.organization_id)
--         OR instance_id IN (NEW.instance_id, ''))
--         AND type = 'organization'
--         AND payload->'organizationScopedUsernames' IS NOT NULL
--     ORDER BY
--         instance_id DESC, organization_id NULLS LAST
--     LIMIT 1;
-- END;
-- $$ LANGUAGE plpgsql;

-- updates the username uniqueness on org settings change
-- CREATE OR REPLACE FUNCTION zitadel.settings_set_username_uniqueness()
-- RETURNS TRIGGER AS $$
-- BEGIN
    -- UPDATE zitadel.users
    -- SET username_org_unique = (NEW.payload->'organizationScopedUsernames')::BOOLEAN
    -- WHERE 
    --     (instance_id = NEW.instance_id AND organization_id = NEW.organization_id)
    --     OR (instance_id = NEW.instance_id AND NEW.organization_id IS NULL)
-- END;
-- $$ LANGUAGE plpgsql;

-- ------------------------------------------------------------
-- triggers
-- ------------------------------------------------------------

-- CREATE TRIGGER  trg_username_uniqueness
-- BEFORE INSERT ON zitadel.users
-- FOR EACH ROW
-- WHEN (NEW.username_org_unique IS NULL)
-- EXECUTE FUNCTION zitadel.user_set_username_uniqueness();

-- CREATE TRIGGER  trg_user_username_uniqueness
-- BEFORE INSERT ON zitadel.human_users
-- FOR EACH ROW
-- WHEN (NEW.username_org_unique IS NULL)
-- EXECUTE FUNCTION zitadel.user_set_username_uniqueness();

-- CREATE TRIGGER  trg_user_username_uniqueness
-- BEFORE INSERT ON zitadel.machine_users
-- FOR EACH ROW
-- WHEN (NEW.username_org_unique IS NULL)
-- EXECUTE FUNCTION zitadel.user_set_username_uniqueness();

CREATE TRIGGER trg_set_updated_at
BEFORE UPDATE ON zitadel.users
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE TRIGGER trg_set_updated_at
BEFORE UPDATE ON zitadel.human_users
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE TRIGGER trg_set_updated_at
BEFORE UPDATE ON zitadel.machine_users
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

-- CREATE TRIGGER  trg_username_uniqueness
-- BEFORE INSERT OR UPDATE ON zitadel.settings
-- FOR EACH ROW
-- WHEN (NEW.type = 'organization' AND OLD.payload->'organizationScopedUsernames' IS DISTINCT FROM NEW.payload->'organizationScopedUsernames')
-- EXECUTE FUNCTION zitadel.org_set_username_uniqueness();
