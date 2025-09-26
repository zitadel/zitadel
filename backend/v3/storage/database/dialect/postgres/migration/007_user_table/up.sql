-- ------------------------------------------------------------
-- table definitions
-- ------------------------------------------------------------

CREATE TYPE zitadel.user_state AS ENUM (
    'active'
    , 'inactive'
    , 'locked'
    , 'suspended'
);

-- user
CREATE TABLE zitadel.users(
    instance_id TEXT NOT NULL
    , org_id TEXT NOT NULL
    , id TEXT NOT NULL


    , username TEXT NOT NULL
    , username_org_unique BOOLEAN NOT NULL -- this field MUST be filled if the username must be unique on organization level
    , state zitadel.user_state NOT NULL

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    
    , PRIMARY KEY (instance_id, org_id, id)
    , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations(instance_id, id)
);

CREATE UNIQUE INDEX ON zitadel.users(instance_id, org_id, username) WHERE username_org_unique IS TRUE;
CREATE UNIQUE INDEX ON zitadel.users(instance_id, username) WHERE username_org_unique IS FALSE;

-- machine user
CREATE TABLE zitadel.machine_users(
    name TEXT
    , description TEXT

    , PRIMARY KEY (instance_id, org_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
    , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations
) INHERITS (zitadel.users);

CREATE INDEX idx_machine_name ON zitadel.machine_users (instance_id, name);
CREATE INDEX idx_machine_user_username ON zitadel.machine_users (instance_id, username);
CREATE INDEX idx_machine_user_username_insensitive ON zitadel.machine_users (instance_id, lower(username));

-- human user
CREATE TABLE zitadel.human_users(
    first_name TEXT
    , last_name TEXT
    , nick_name TEXT
    , display_name TEXT
    , preferred_language TEXT
    , gender SMALLINT 
    , avatar_key TEXT

    , PRIMARY KEY (instance_id, org_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
    , FOREIGN KEY (instance_id, org_id) REFERENCES zitadel.organizations
) INHERITS (zitadel.users);

CREATE INDEX idx_human_user_username ON zitadel.human_users (instance_id, username);
CREATE INDEX idx_human_user_username_insensitive ON zitadel.human_users (instance_id, lower(username));

CREATE TYPE zitadel.human_contact_type AS ENUM (
    'email'
    , 'phone'
);

CREATE TABLE zitadel.human_contacts(
    instance_id TEXT NOT NULL
    , org_id TEXT NOT NULL
    , user_id TEXT NOT NULL
    , type human_contact_type NOT NULL
    , current_value TEXT
    , is_verified BOOLEAN NOT NULL DEFAULT FALSE
    , unverified_value TEXT -- if a user wants to update the info but its not yet verified, verification is done in a separate issue

    , PRIMARY KEY (instance_id, org_id, user_id, type)
    , FOREIGN KEY (instance_id, org_id, user_id) REFERENCES zitadel.human_users(instance_id, org_id, id) ON DELETE CASCADE
);

CREATE INDEX idx_human_contacts_value ON zitadel.human_contacts(instance_id, current_value);
CREATE INDEX idx_human_contacts_value_lower ON zitadel.human_contacts(instance_id, lower(current_value)) WHERE type = 'email';

-- ------------------------------------------------------------
-- function definitions
-- ------------------------------------------------------------

-- sets the username uniqueness initially
CREATE OR REPLACE FUNCTION zitadel.user_set_username_uniqueness()
RETURNS TRIGGER AS $$
BEGIN
    SELECT 
        payload->'organizationScopedUsernames'::BOOLEAN INTO NEW.username_org_unique 
    FROM 
        zitadel.settings 
    WHERE 
        ((instance_id = NEW.instance_id AND org_id = NEW.org_id)
        OR instance_id IN (NEW.instance_id, ''))
        AND type = 'organization'
        AND payload->'organizationScopedUsernames' IS NOT NULL;
    ORDER BY
        instance_id DESC, org_id NULLS LAST
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- updates the username uniqueness on org settings change
CREATE OR REPLACE FUNCTION zitadel.settings_set_username_uniqueness()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE zitadel.users
    SET username_org_unique = (NEW.payload->'organizationScopedUsernames')::BOOLEAN
    WHERE 
        (instance_id = NEW.instance_id AND org_id = NEW.org_id)
        OR (instance_id = NEW.instance_id AND NEW.org_id IS NULL)
END;
$$ LANGUAGE plpgsql;

-- ------------------------------------------------------------
-- triggers
-- ------------------------------------------------------------

CREATE TRIGGER  trg_username_uniqueness
BEFORE INSERT ON zitadel.users
FOR EACH ROW
WHEN (NEW.username_org_unique IS NULL)
EXECUTE FUNCTION zitadel.user_set_username_uniqueness();

CREATE TRIGGER  trg_user_username_uniqueness
BEFORE INSERT ON zitadel.human_users
FOR EACH ROW
WHEN (NEW.username_org_unique IS NULL)
EXECUTE FUNCTION zitadel.user_set_username_uniqueness();

CREATE TRIGGER  trg_user_username_uniqueness
BEFORE INSERT ON zitadel.machine_users
FOR EACH ROW
WHEN (NEW.username_org_unique IS NULL)
EXECUTE FUNCTION zitadel.user_set_username_uniqueness();

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

CREATE TRIGGER  trg_username_uniqueness
BEFORE INSERT OR UPDATE ON zitadel.settings
FOR EACH ROW
WHEN (NEW.type = 'organization' AND OLD.payload->'organizationScopedUsernames' IS DISTINCT FROM NEW.payload->'organizationScopedUsernames')
EXECUTE FUNCTION zitadel.org_set_username_uniqueness();
