CREATE SCHEMA adminapi;

CREATE TABLE adminapi.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, instance_id)
);

CREATE TABLE adminapi.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    event_date TIMESTAMPTZ,
    last_successful_spooler_run TIMESTAMPTZ,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, instance_id)
);

CREATE TABLE adminapi.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (view_name, failed_sequence, instance_id)
);

CREATE TABLE adminapi.styling (
    aggregate_id TEXT NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    label_policy_state INT2 NOT NULL DEFAULT 0::INT2,
    sequence INT8 NULL,
    primary_color TEXT NULL,
    background_color TEXT NULL,
    warn_color TEXT NULL,
    font_color TEXT NULL,
    primary_color_dark TEXT NULL,
    background_color_dark TEXT NULL,
    warn_color_dark TEXT NULL,
    font_color_dark TEXT NULL,
    logo_url TEXT NULL,
    icon_url TEXT NULL,
    logo_dark_url TEXT NULL,
    icon_dark_url TEXT NULL,
    font_url TEXT NULL,
    err_msg_popup BOOL NULL,
    disable_watermark BOOL NULL,
    hide_login_name_suffix BOOL NULL,
    instance_id TEXT NOT NULL,

    PRIMARY KEY (aggregate_id, label_policy_state, instance_id)
);
