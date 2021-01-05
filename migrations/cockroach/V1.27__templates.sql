CREATE TABLE adminapi.mail_templates (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    mail_template_state SMALLINT,
    sequence BIGINT,

    template BYTES,

    PRIMARY KEY (aggregate_id)
);


CREATE TABLE management.mail_templates (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    mail_template_state SMALLINT,
    sequence BIGINT,

    template BYTES,

    PRIMARY KEY (aggregate_id)
);

GRANT SELECT ON TABLE adminapi.mail_templates TO notification;
GRANT SELECT ON TABLE management.mail_templates TO notification;


CREATE TABLE adminapi.mail_texts (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    mail_text_state SMALLINT,
    sequence BIGINT,

    mail_text_type TEXT,
    language TEXT,
    title TEXT,
    pre_header TEXT,
    subject TEXT,
    greeting TEXT,
    text TEXT,
    button_text TEXT,

    PRIMARY KEY (aggregate_id, mail_text_type, language)
);


CREATE TABLE management.mail_texts (
    aggregate_id TEXT,

    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    mail_text_state SMALLINT,
    sequence BIGINT,

    mail_text_type TEXT,
    language TEXT,
    title TEXT,
    pre_header TEXT,
    subject TEXT,
    greeting TEXT,
    text TEXT,
    button_text TEXT,

    PRIMARY KEY (aggregate_id, mail_text_type, language)
);

GRANT SELECT ON TABLE adminapi.mail_texts TO notification;
GRANT SELECT ON TABLE management.mail_texts TO notification;


ALTER TABLE management.project_roles ADD COLUMN change_date TIMESTAMPTZ;
ALTER TABLE auth.project_roles ADD COLUMN change_date TIMESTAMPTZ;