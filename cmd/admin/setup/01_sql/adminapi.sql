CREATE SCHEMA adminapi;

CREATE TABLE adminapi.locks (
    locker_id TEXT,
    locked_until TIMESTAMPTZ(3),
    projection_name TEXT,

    PRIMARY KEY (projection_name)
);

CREATE TABLE adminapi.current_sequences (
    projection_name TEXT,
    aggregate_type TEXT,
    current_sequence BIGINT,
    timestamp TIMESTAMPTZ,

    PRIMARY KEY (projection_name, aggregate_type)
);

CREATE TABLE adminapi.failed_events (
    projection_name TEXT,
    failed_sequence BIGINT,
    failure_count SMALLINT,
    error TEXT,

    PRIMARY KEY (projection_name, failed_sequence)
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

    PRIMARY KEY (aggregate_id, label_policy_state)
);
