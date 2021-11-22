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
