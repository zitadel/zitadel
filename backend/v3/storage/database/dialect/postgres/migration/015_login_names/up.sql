CREATE TABLE zitadel.login_names(
	instance_id TEXT NOT NULL
    , organization_id TEXT -- is set if the setting for loginNameIncludesDomain is set on organization level, otherwise null
	, user_id TEXT NOT NULL
	
    , username TEXT NOT NULL CHECK (username <> '')
    , domain TEXT
    , login_name TEXT GENERATED ALWAYS AS (CASE WHEN domain IS NOT NULL THEN username || '@' || domain ELSE username END) STORED
	, is_preferred BOOLEAN NOT NULL DEFAULT FALSE

    , used_setting TEXT NOT NULL

	, PRIMARY KEY (instance_id, user_id, login_name)
	, FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, organization_id, domain) REFERENCES zitadel.org_domains(instance_id, org_id, domain) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_login_names_instance_login_name ON zitadel.login_names(instance_id, lower(login_name));
CREATE INDEX idx_login_names_instance_user ON zitadel.login_names(instance_id, user_id);
CREATE INDEX idx_login_names_setting ON zitadel.login_names(instance_id, used_setting);
CREATE INDEX idx_login_names_domain ON zitadel.login_names(instance_id, domain); -- used for cleanup of login names when a domain is deleted

CREATE OR REPLACE FUNCTION zitadel.apply_domain_manipulation_to_login_names() RETURNS TRIGGER AS $$
DECLARE
    setting zitadel.settings%ROWTYPE;
BEGIN
    IF (NOT NEW.is_verified) THEN
        -- TODO(adlerhurst): is it possible that the domain is updated from verified to unverified?
        RAISE NOTICE 'Domain is not verified, skipping login name manipulation';
        RETURN NULL;
    END IF;

    SELECT settings.*
    INTO setting
    FROM zitadel.settings
    WHERE
        settings.instance_id = NEW.instance_id
        AND settings.type = 'domain'
        AND (settings.organization_id IS NULL OR settings.organization_id = NEW.org_id)
    ORDER BY settings.organization_id NULLS LAST
    LIMIT 1;

    IF NOT (setting.settings->'loginNameIncludesDomain')::BOOLEAN THEN
        RAISE NOTICE 'skip adding domain because loginNameIncludesDomain is false for setting';
        RETURN NULL;
    END IF;

    -- Lock based on the setting scope (organization-specific or instance-level/global).
    -- The lock is released at transaction end.
    PERFORM pg_advisory_xact_lock(
        hashtext('zitadel.login_names')
        , hashtext(setting.instance_id || ':' || COALESCE(setting.organization_id, 'global'))
    );

    IF NEW.is_primary IS DISTINCT FROM OLD.is_primary THEN
        RAISE NOTICE 'primary domain changed, updating preferred login name';
        UPDATE zitadel.login_names
        SET is_preferred = NEW.is_primary
        WHERE
            login_names.instance_id = NEW.instance_id
            AND login_names.domain = NEW.domain;
    END IF;

    IF NEW.domain IS DISTINCT FROM OLD.domain THEN
        RAISE NOTICE 'domain changed, updating login names';
        UPDATE zitadel.login_names
        SET domain = NEW.domain
        WHERE
            login_names.instance_id = NEW.instance_id
            AND login_names.organization_id = NEW.org_id
            AND login_names.domain = OLD.domain;
    END IF;

    IF NEW.is_verified IS DISTINCT FROM OLD.is_verified THEN
        RAISE NOTICE 'Domain verification changed, inserting login names for verified domain';
        INSERT INTO zitadel.login_names(instance_id, organization_id, user_id, username, domain, is_preferred, used_setting)
        SELECT
            NEW.instance_id
            , NEW.org_id
            , users.id
            , users.username
            , NEW.domain
            , NEW.is_primary
            , setting.id
        FROM
            zitadel.users
        WHERE
            users.instance_id = NEW.instance_id
            AND users.organization_id = NEW.org_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_apply_domain_manipulation_to_login_names
AFTER INSERT OR UPDATE ON zitadel.org_domains
FOR EACH ROW
EXECUTE FUNCTION zitadel.apply_domain_manipulation_to_login_names();

CREATE OR REPLACE FUNCTION zitadel.apply_user_update_to_login_names() RETURNS TRIGGER AS $$
BEGIN
    -- Lock global scope and organization scope to serialize with instance-level
    -- and org-level setting/domain manipulations.
    -- The lock is released at transaction end.
    PERFORM pg_advisory_xact_lock(hashtext('zitadel.login_names'), hashtext(NEW.instance_id || ':global'));
    PERFORM pg_advisory_xact_lock(hashtext('zitadel.login_names'), hashtext(NEW.instance_id || ':' || NEW.organization_id));

    UPDATE
        zitadel.login_names
    SET
        username = NEW.username
    WHERE
        login_names.instance_id = NEW.instance_id
        AND login_names.user_id = NEW.id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_apply_user_update_to_login_names
AFTER UPDATE ON zitadel.users
FOR EACH ROW
WHEN (NEW.username IS DISTINCT FROM OLD.username)
EXECUTE FUNCTION zitadel.apply_user_update_to_login_names();

CREATE OR REPLACE FUNCTION zitadel.apply_user_insert_to_login_names() RETURNS TRIGGER AS $$
DECLARE
    setting zitadel.settings%ROWTYPE;
BEGIN
    SELECT settings.*
    INTO setting
    FROM
        zitadel.settings
    WHERE
        settings.instance_id = NEW.instance_id
        AND settings.type = 'domain'
        AND (
            settings.organization_id IS NULL
            OR settings.organization_id = NEW.organization_id
        )
    ORDER BY settings.organization_id NULLS LAST
    LIMIT 1;

    -- Lock based on the setting scope (organization-specific or instance-level/global).
    -- The lock is released at transaction end.
    PERFORM pg_advisory_xact_lock(
        hashtext('zitadel.login_names')
        , hashtext(setting.instance_id || ':' || COALESCE(setting.organization_id, 'global'))
    );

    IF NOT (setting.settings->'loginNameIncludesDomain')::BOOLEAN THEN
        RAISE NOTICE 'inserting username as login name';
        INSERT INTO zitadel.login_names(instance_id, organization_id, user_id, username, is_preferred, used_setting)
        VALUES (NEW.instance_id, NEW.organization_id, NEW.id, NEW.username, TRUE, setting.id);
        
        RETURN NULL;
    END IF;

    RAISE NOTICE 'inserting login names';
    INSERT INTO zitadel.login_names(instance_id, organization_id, user_id, username, domain, is_preferred, used_setting)
    SELECT
        NEW.instance_id
        , NEW.organization_id
        , NEW.id
        , NEW.username
        , org_domains.domain
        , org_domains.is_primary IS NULL OR org_domains.is_primary
        , setting.id
    FROM
        zitadel.org_domains
    WHERE
        org_domains.instance_id = NEW.instance_id
        AND org_domains.org_id = NEW.organization_id
        AND org_domains.is_verified = TRUE;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_apply_user_insert_to_login_names
AFTER INSERT ON zitadel.users
FOR EACH ROW
EXECUTE FUNCTION zitadel.apply_user_insert_to_login_names();

CREATE OR REPLACE FUNCTION zitadel.apply_domain_policy_manipulation_to_login_names() RETURNS TRIGGER AS $$
DECLARE
    old_lock_key INT4;
    new_lock_key INT4;
BEGIN
    CASE TG_OP
        -- insert and delete can only happen on organization level therefore we can load the instance level setting as OLD or NEW.
        WHEN 'DELETE' THEN
            SELECT
                settings.*
            INTO
                NEW
            FROM
                zitadel.settings
            WHERE
                settings.instance_id = OLD.instance_id
                AND settings.type = 'domain'
                AND settings.organization_id IS NULL;
        WHEN 'INSERT' THEN
            SELECT
                settings.*
            INTO
                OLD
            FROM
                zitadel.settings
            WHERE
                settings.instance_id = NEW.instance_id
                AND settings.type = 'domain'
                AND settings.organization_id IS NULL;
        ELSE
            -- OLD and NEW are already populated, do nothing
    END CASE;

    -- Lock both OLD and NEW setting scopes to serialize transitions across settings.
    -- Acquire in deterministic order to avoid deadlocks.
    old_lock_key := hashtext(OLD.instance_id || ':' || COALESCE(OLD.organization_id, 'global'));
    new_lock_key := hashtext(NEW.instance_id || ':' || COALESCE(NEW.organization_id, 'global'));
    IF old_lock_key = new_lock_key THEN
        PERFORM pg_advisory_xact_lock(hashtext('zitadel.login_names'), old_lock_key);
    ELSE
        PERFORM pg_advisory_xact_lock(hashtext('zitadel.login_names'), new_lock_key);
        PERFORM pg_advisory_xact_lock(hashtext('zitadel.login_names'), old_lock_key);
    END IF;
    
    IF (OLD.settings->'loginNameIncludesDomain')::BOOLEAN IS NOT DISTINCT FROM (NEW.settings->'loginNameIncludesDomain')::BOOLEAN THEN
        -- field not changed but the setting id did change.
        IF TG_OP = 'DELETE' OR TG_OP = 'INSERT' THEN
            UPDATE
                zitadel.login_names
            SET
                used_setting = NEW.id
            WHERE
                login_names.instance_id = OLD.instance_id
                AND login_names.used_setting = OLD.id
                AND (NEW.organization_id IS NULL OR login_names.organization_id = NEW.organization_id);
        END IF;

        -- no change in the setting, so we can skip the update
        RETURN NULL;
    END IF;

    RAISE NOTICE 'recompute login names';

    WITH affected_users AS (
        DELETE FROM zitadel.login_names
        WHERE
            login_names.instance_id = NEW.instance_id
            AND login_names.used_setting = OLD.id
            AND (NEW.organization_id IS NULL OR login_names.organization_id = NEW.organization_id)
        RETURNING instance_id, organization_id, user_id
    )
    INSERT INTO zitadel.login_names(instance_id, organization_id, user_id, username, domain, is_preferred, used_setting)
    SELECT DISTINCT ON (au.instance_id, au.user_id, users.username, org_domains.domain)
        au.instance_id
        , au.organization_id
        , au.user_id
        , users.username
        , org_domains.domain
        , org_domains.is_primary IS NULL OR org_domains.is_primary
        , NEW.id
    FROM
        affected_users au
    JOIN zitadel.users ON
        users.instance_id = au.instance_id
        AND users.id = au.user_id
    LEFT JOIN zitadel.org_domains ON 
        org_domains.instance_id = au.instance_id
        AND org_domains.org_id = au.organization_id
        AND org_domains.is_verified
        AND (NEW.settings->'loginNameIncludesDomain')::BOOLEAN;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_apply_domain_policy_manipulation_to_login_names
AFTER INSERT OR UPDATE ON zitadel.settings
FOR EACH ROW
WHEN (NEW.type = 'domain')
EXECUTE FUNCTION zitadel.apply_domain_policy_manipulation_to_login_names();

CREATE OR REPLACE TRIGGER trg_apply_domain_policy_removed_to_login_names
AFTER DELETE ON zitadel.settings
FOR EACH ROW
WHEN (OLD.type = 'domain')
EXECUTE FUNCTION zitadel.apply_domain_policy_manipulation_to_login_names();