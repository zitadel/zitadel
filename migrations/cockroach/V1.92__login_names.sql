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

   
-- --------------------------------------------------------
-- only default domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    ('org-ch.localhost', true, 'org') 
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- default and additional domain verified
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    ('org-ch.localhost', true, 'org') 
    , ('custom.ch', false, 'org')
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

-- default and custom login name => 4 login names default is primary
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- default and additional domain verified and primary
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    ('org-ch.localhost', false, 'org') 
    , ('custom.ch', true, 'org')
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

-- default and custom login name => 2 login names default is primary
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) no domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    -- only default for org
    ('org-ch.localhost', true, 'org') 
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- default and custom login name => 1 login name machine=user_name human=user_name(=email)
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) verified domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', true, 'org') 
    , ('custom.ch', false, 'org')
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- 1 login name machine=user_name human=user_name(=email)
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) verified, primary domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', false, 'org') 
    , ('custom.ch', true, 'org')
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- 1 login name machine=user_name human=user_name(=email)
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) no domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    -- only default for org
    ('org-ch.localhost', true, 'org') 
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- one login per user
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) verified domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', true, 'org') 
    , ('custom.ch', false, 'org')
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- 2 login names per user
SELECT * FROM zitadel.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) verified, primary domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO zitadel.projections.login_names_users (
    id
    , user_name
    , resource_owner
) VALUES 
    ('h', 'human', 'org')
    , ('m', 'machine', 'org')
;

INSERT INTO zitadel.projections.login_names_domains (
    name
    , is_primary
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', false, 'org') 
    , ('custom.ch', true, 'org')
;

INSERT INTO zitadel.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- 2 login names per user
explain analyze SELECT * FROM zitadel.projections.login_names WHERE user_id = 'h';
ROLLBACK;
