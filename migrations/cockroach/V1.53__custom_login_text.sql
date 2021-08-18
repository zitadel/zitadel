CREATE TABLE adminapi.custom_texts (
     aggregate_id TEXT,

     creation_date TIMESTAMPTZ,
     change_date TIMESTAMPTZ,
     sequence BIGINT,

     template TEXT,
     language TEXT,
     key TEXT,
     text TEXT,

     PRIMARY KEY (aggregate_id, template, key, language)
);

CREATE TABLE management.custom_texts (
      aggregate_id TEXT,

      creation_date TIMESTAMPTZ,
      change_date TIMESTAMPTZ,
      sequence BIGINT,

      template TEXT,
      language TEXT,
      key TEXT,
      text TEXT,

      PRIMARY KEY (aggregate_id, template, key, language)
);

CREATE TABLE auth.custom_texts (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    template TEXT,
    language TEXT,
    key TEXT,
    text TEXT,

    PRIMARY KEY (aggregate_id, template, key, language)
);
