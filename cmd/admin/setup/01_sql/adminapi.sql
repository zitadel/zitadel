CREATE SCHEMA adminapi;

CREATE TABLE adminapi.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    view_name TEXT,

    PRIMARY KEY (view_name)
);

CREATE TABLE adminapi.current_sequences (
    view_name TEXT,
    current_sequence BIGINT,
    event_timestamp TIMESTAMPTZ,
    last_successful_spooler_run TIMESTAMPTZ,

    PRIMARY KEY (view_name)
);

CREATE TABLE adminapi.failed_events (
    view_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    err_msg TEXT,

    PRIMARY KEY (view_name, failed_sequence)
);

CREATE TABLE adminapi.styling (
    aggregate_id STRING NOT NULL,
    creation_date TIMESTAMPTZ NULL,
    change_date TIMESTAMPTZ NULL,
    label_policy_state INT2 NOT NULL DEFAULT 0:::INT2,
    sequence INT8 NULL,
    primary_color STRING NULL,
    background_color STRING NULL,
    warn_color STRING NULL,
    font_color STRING NULL,
    primary_color_dark STRING NULL,
    background_color_dark STRING NULL,
    warn_color_dark STRING NULL,
    font_color_dark STRING NULL,
    logo_url STRING NULL,
    icon_url STRING NULL,
    logo_dark_url STRING NULL,
    icon_dark_url STRING NULL,
    font_url STRING NULL,
    err_msg_popup BOOL NULL,
    disable_watermark BOOL NULL,
    hide_login_name_suffix BOOL NULL,
    instance_id STRING NOT NULL,

    PRIMARY KEY (aggregate_id, label_policy_state)
);
