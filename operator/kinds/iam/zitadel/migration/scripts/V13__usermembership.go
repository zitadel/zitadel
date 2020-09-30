package scripts

const V13Usermembership = `
CREATE TABLE management.user_memberships (
    user_id TEXT,
    member_type SMALLINT,
    aggregate_id TEXT,
    object_id TEXT,

    roles TEXT ARRAY,
    display_name TEXT,
    resource_owner TEXT,
    creation_date TIMESTAMPTZ,
    change_date TIMESTAMPTZ,
    sequence BIGINT,

    PRIMARY KEY (user_id, member_type, aggregate_id, object_id)
);
`
