CREATE TABLE test.projections.login_names_users (
    id STRING NOT NULL,
    type SMALLINT NOT NULL,
    user_name STRING NOT NULL,
    email STRING,
    is_domain_claimed BOOLEAN NOT NULL DEFAULT false,
    resource_owner STRING NOT NULL,

    PRIMARY KEY (id),
    INDEX idx_ro (resource_owner)
);

CREATE TABLE test.projections.login_names_domains (
    name STRING NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    resource_owner STRING NOT NULL,

    PRIMARY KEY (resource_owner, name)
);

CREATE TABLE test.projections.login_names_policies (
    must_be_domain BOOLEAN NOT NULL,
    is_default BOOLEAN NOT NULL,
    resource_owner STRING NOT NULL,

    PRIMARY KEY (resource_owner)
);


-- drop VIEW test.projections.login_names;
CREATE VIEW test.projections.login_names
AS SELECT
    user_id
    , type AS user_type
    , IF(
        must_be_domain
        , CONCAT(user_name, '@', domain)
        , CASE type
            WHEN 1 THEN email --human
            WHEN 2 THEN user_name --machine
        END
    ) AS login_name
    , IFNULL(is_primary, true) AS is_primary -- is_default is null no additional verified domain and policy with must_be_domain=false
FROM (
SELECT
    p.user_id
    , p.type
    , p.user_name
    , p.email
    , p.resource_owner
    , p.must_be_domain
    , domains.name AS domain
    , domains.is_primary 
FROM (
    SELECT 
        users.id as user_id
        , users.type
        , users.user_name
        , users.email
        , users.resource_owner
        , IFNULL(policy_custom.must_be_domain, policy_default.must_be_domain) must_be_domain
    FROM test.projections.login_names_users users
    LEFT JOIN test.projections.login_names_policies policy_custom on policy_custom.resource_owner = users.resource_owner
    LEFT JOIN test.projections.login_names_policies policy_default on policy_default.is_default = true) p
LEFT JOIN test.projections.login_names_domains domains ON p.must_be_domain AND domains.is_verified AND p.resource_owner = domains.resource_owner
);

-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- domain
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------

-- --------------------------------------------------------
-- only default domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    ('org-ch.localhost', true, true, 'org') 
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- default and additional domain unverified
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    ('org-ch.localhost', true, true, 'org') 
    , ('custom.ch', false, false, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

-- only default login name, because second domain not verified => 1 login names
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- default and additional domain verified
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    ('org-ch.localhost', true, true, 'org') 
    , ('custom.ch', false, true, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

-- default and custom login name => 2 login names default is primary
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- default and additional domain verified and primary
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    ('org-ch.localhost', false, true, 'org') 
    , ('custom.ch', true, true, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
;

-- default and custom login name => 2 login names default is primary
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- custom policy (must_be_domain=false)
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) no domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- only default for org
    ('org-ch.localhost', true, true, 'org') 
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- default and custom login name => 1 login name machine=user_name human=email
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) unverified domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', true, true, 'org') 
    , ('custom.ch', false, false, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- login name machine=user_name human=email
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) verified domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', true, true, 'org') 
    , ('custom.ch', false, true, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- 1 login name machine=user_name human=email
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=false) verified, primary domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', false, true, 'org') 
    , ('custom.ch', true, true, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (false, false, 'org')
;

-- 1 login name machine=user_name human=email
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- custom policy (must_be_domain=true)
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) no domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- only default for org
    ('org-ch.localhost', true, true, 'org') 
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- one login per user
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) unverified domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', true, true, 'org') 
    , ('custom.ch', false, false, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- one login per user
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) verified domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', true, true, 'org') 
    , ('custom.ch', false, true, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- 2 login names per user
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- custom policy (must_be_domain=true) verified, primary domain
-- --------------------------------------------------------
BEGIN;

INSERT INTO test.projections.login_names_users (
    id
    , type
    , user_name
    , email
    , is_domain_claimed
    , resource_owner
) VALUES 
    ('h', 1, 'human', 'human@caos.ch', false, 'org')
    , ('m', 2, 'machine', NULL, false, 'org')
;

INSERT INTO test.projections.login_names_domains (
    name
    , is_primary
    , is_verified
    , resource_owner
) VALUES
    -- default and unverified for org
    ('org-ch.localhost', false, true, 'org') 
    , ('custom.ch', true, true, 'org')
;

INSERT INTO test.projections.login_names_policies (
    must_be_domain
    , is_default
    , resource_owner
) VALUES 
    (true, true, 'IAM')
    , (true, false, 'org')
;

-- 2 login names per user
SELECT * FROM test.projections.login_names WHERE user_id IN ('h', 'm');
ROLLBACK;

-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- domain claimed: correct implementation is impossible atm
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
-- --------------------------------------------------------
