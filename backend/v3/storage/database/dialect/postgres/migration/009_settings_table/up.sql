CREATE TYPE zitadel.settings_type AS ENUM (
    'login',
    'branding',
    'password_complexity',
    'password_expiry',
    'domain',
    'lockout',
    'security',
    'organization',
    'notification',
    'legal_and_support',
    'secret_generator'
);

CREATE TYPE zitadel.settings_state AS ENUM (
    'active',
    'preview'
);

CREATE TABLE zitadel.settings (
    instance_id TEXT NOT NULL
    , organization_id TEXT
    , id TEXT NOT NULL DEFAULT gen_random_uuid()
    , type zitadel.settings_type NOT NULL
    , state zitadel.settings_state NOT NULL
    , attributes JSONB -- the storage does not really care about what is configured so we store it as json

    , created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    , updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    , PRIMARY KEY (instance_id, id, type, state)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, organization_id) REFERENCES zitadel.organizations(instance_id, id) ON DELETE CASCADE

    , UNIQUE (instance_id, organization_id, type, state)
);

CREATE UNIQUE INDEX idx_settings_unique_type ON zitadel.settings (instance_id, organization_id, type, state) NULLS NOT DISTINCT;

CREATE TRIGGER trigger_set_updated_at
BEFORE UPDATE ON zitadel.settings
FOR EACH ROW
WHEN (NEW.updated_at IS NULL)
EXECUTE FUNCTION zitadel.set_updated_at();

CREATE OR REPLACE FUNCTION zitadel.jsonb_patch(
  INOUT source JSONB
  , path TEXT[]
  , p_value ANYELEMENT
  , is_array BOOLEAN DEFAULT FALSE -- indicates if the property is an array. If true p_value is added to the array or removed if delete_array_element is true.
  , delete_array_element BOOLEAN DEFAULT FALSE -- if true, p_value is removed from the array instead of added (only if is_array is true)
)
IMMUTABLE
PARALLEL SAFE
COST 5
LANGUAGE 'plpgsql'
AS $$
  BEGIN
    IF source #> path[1:array_length(path, 1)-1] IS NULL THEN
      source := zitadel.jsonb_patch(source, path[1:array_length(path, 1)-1], '{}'::JSONB);
    END IF;

    IF is_array THEN
      IF delete_array_element THEN
        source := jsonb_set(source, path, (SELECT jsonb_agg(elem) FROM jsonb_array_elements(source #> path) AS elem WHERE elem <> to_jsonb(p_value)));
        RETURN;
      END IF;
      IF source #> path IS NULL THEN
        source := jsonb_set(source, path, jsonb_build_array(p_value));
        RETURN;
      END IF;

      source := jsonb_set(source, path, source #> path || to_jsonb(p_value));
      RETURN;
    END IF;

    IF p_value IS NULL THEN
      source := jsonb_set_lax(source, path, NULL, true, 'delete_key');
      RETURN;
    END IF;
    source := jsonb_set_lax(source, path, to_jsonb(p_value), true, 'delete_key');
  END;
$$;
