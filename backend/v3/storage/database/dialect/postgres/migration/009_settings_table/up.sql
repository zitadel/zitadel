CREATE TYPE zitadel.settings_type AS ENUM (
    'login',
    'branding',
    'password_complexity',
    'password_expiry',
    'domain',
    'lockout',
    'security',
    'organization'
    'notification'
    'legal_and_support'
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
    , settings JSONB -- the storage does not really care about what is configured so we store it as json

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

CREATE OR REPLACE FUNCTION zitadel.jsonb_array_remove(
    source JSONB
    , path TEXT[]
    , value anyelement

    , outt out jsonb
)
    language 'plpgsql'
AS $$
  BEGIN
    outt = jsonb_set(source, path,
    (CASE WHEN (SELECT source ?& path) then
      coalesce(
      (SELECT jsonb_agg(v)
      FROM jsonb_array_elements(source #>path) AS elem(v)
      WHERE v::text <> value::text),
      jsonb_build_array())
    ELSe
      jsonb_build_array()
    END)::jsonb, TRUE);
  END;
$$;

CREATE OR REPLACE FUNCTION zitadel.jsonb_array_append(
    source jsonb
    , path text[]
    , value anyelement

    , outt out jsonb
)
    language 'plpgsql'
AS $$
  BEGIN
    outt := jsonb_insert(
      source,
      (CASE WHEN (SELECT source ?& path) then
        array_append(path, '-1')
      ELSE
        path
      END)::text[],
      (CASE WHEN (select source ?& path) then
        value::TEXT::jsonb
      ELSE
        jsonb_build_array(value)
      END),
      TRUE
  );

  END;
$$;

