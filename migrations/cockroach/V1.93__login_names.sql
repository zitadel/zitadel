CREATE TABLE zitadel.projections.login_names_users (
    id STRING NOT NULL
    , user_name STRING NOT NULL
    , resource_owner STRING NOT NULL

    , PRIMARY KEY (id)
    , INDEX idx_ro (resource_owner)
);

CREATE TABLE zitadel.projections.login_names_domains (
    name STRING NOT NULL
    , is_primary BOOLEAN NOT NULL DEFAULT false
    , resource_owner STRING NOT NULL
    
    , PRIMARY KEY (resource_owner, name)
);

CREATE TABLE zitadel.projections.login_names_policies (
    must_be_domain BOOLEAN NOT NULL
    , is_default BOOLEAN NOT NULL
    , resource_owner STRING NOT NULL
    
    , PRIMARY KEY (resource_owner)
    , INDEX idx_is_default (resource_owner, is_default)
);

CREATE VIEW zitadel.projections.login_names
AS SELECT
    user_id
    , IF(
        must_be_domain
        , CONCAT(user_name, '@', domain)
        , user_name
    ) AS login_name
    , IFNULL(is_primary, true) AS is_primary -- is_default is null on additional verified domain and policy with must_be_domain=false
FROM (
SELECT
    policy_users.user_id
    , policy_users.user_name
    , policy_users.resource_owner
    , policy_users.must_be_domain
    , domains.name AS domain
    , domains.is_primary 
FROM (
    SELECT 
        users.id as user_id
        , users.user_name
        , users.resource_owner
        , IFNULL(policy_custom.must_be_domain, policy_default.must_be_domain) AS must_be_domain
    FROM zitadel.projections.login_names_users users
    LEFT JOIN zitadel.projections.login_names_policies policy_custom on policy_custom.resource_owner = users.resource_owner
    LEFT JOIN zitadel.projections.login_names_policies policy_default on policy_default.is_default = true) policy_users
LEFT JOIN zitadel.projections.login_names_domains domains ON policy_users.must_be_domain AND policy_users.resource_owner = domains.resource_owner
);
