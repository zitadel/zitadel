INSERT INTO aether.instances (id) VALUES ('instance1');

INSERT INTO aether.collections (instance_id, parent_id, id, type) VALUES 
    ('instance1', NULL, 'org1', 'organization')
    , ('instance1', NULL, 'org2', 'organization')
    , ('instance1', 'org1', 'org1_1', 'organization')
    , ('instance1', 'org1_1', 'org1_1_1', 'organization')
    , ('instance1', 'org1_1_1', 'org1_1_1_1', 'organization')
;

INSERT INTO aether.identities (instance_id, collection_id, id) VALUES 
    ('instance1', 'org1', 'user1')
    , ('instance1', 'org1', 'user2')
    , ('instance1', 'org2', 'user3')
    , ('instance1', NULL, 'user4')
;

-- in this example profile.email is defined on collection level, username is global
INSERT INTO aether.identity_properties (instance_id, collection_id, identity_id, path, value, is_identifier) VALUES 
    ('instance1', 'org1', 'user1', 'profile.email', to_jsonb('user1@example.com'::TEXT), TRUE)
    , ('instance1', NULL, 'user1', 'username', to_jsonb('user1name'::TEXT), TRUE)
    , ('instance1', NULL, 'user1', 'id', to_jsonb(1::NUMERIC), TRUE)
    , ('instance1', 'org1', 'user1', 'profile.phone', to_jsonb('+1234567890'::TEXT), FALSE)
    , ('instance1', 'org1', 'user2', 'profile.email', to_jsonb('user2@example.com'::TEXT), TRUE)
    , ('instance1', NULL, 'user2', 'username', to_jsonb('user2name'::TEXT), TRUE)
    , ('instance1', NULL, 'user2', 'id', to_jsonb(2::NUMERIC), TRUE)
    , ('instance1', 'org2', 'user3', 'profile.email', to_jsonb('user3@example.com'::TEXT), TRUE)
    , ('instance1', NULL, 'user3', 'username', to_jsonb('user3name'::TEXT), TRUE)
    , ('instance1', NULL, 'user3', 'id', to_jsonb(3::NUMERIC), TRUE)
    , ('instance1', NULL, 'user4', 'username', to_jsonb('user4name'::TEXT), TRUE)
    , ('instance1', NULL, 'user4', 'id', to_jsonb(4::NUMERIC), TRUE)
;

SELECT * FROM aether.identity_identifiers;

INSERT INTO aether.settings (instance_id, collection_id, type) VALUES 
    ('instance1', NULL, 'notification')
    , ('instance1', 'org1', 'notification')
    , ('instance1', 'org1_1_1_1', 'notification')
;

INSERT INTO aether.setting_properties (instance_id, setting_id, path, value) VALUES 
    -- Global notification settings
    ('instance1', 1, 'email.enabled', to_jsonb(true::BOOLEAN))
    , ('instance1', 2, 'email.enabled', to_jsonb(false::BOOLEAN))
    , ('instance1', 2, 'sms.enabled', to_jsonb(true::BOOLEAN))
    , ('instance1', 3, 'push.enabled', to_jsonb(true::BOOLEAN))
;

-- try project as collection
INSERT INTO aether.collections (instance_id, parent_id, id, type) VALUES 
    ('instance1', NULL, 'Zitadel', 'project')
    
    , ('instance1', 'Zitadel', 'roles', 'roles')
    , ('instance1', 'roles', 'admin', 'role')
    , ('instance1', 'roles', 'developer', 'role')
    , ('instance1', 'roles', 'user', 'role')

    , ('instance1', 'Zitadel', 'grants', 'grants')
    , ('instance1', 'grants', 'grant-org1', 'grant')
    , ('instance1', 'grants', 'grant-org2', 'grant')

    -- , ('instance1', 'Zitadel', 'authorizations')
    -- , ('instance1', 'Zitadel', 'Console')
;

INSERT INTO aether.collection_properties (instance_id, collection_id, path, value, is_linking) VALUES 
    ('instance1', 'Zitadel', 'project.id', to_jsonb('zitadel-project'::TEXT), FALSE)
    , ('instance1', 'Zitadel', 'project.name', to_jsonb('Zitadel Project'::TEXT), FALSE)

    , ('instance1', 'roles', 'role.type', to_jsonb('base'::TEXT), FALSE)

    , ('instance1', 'admin', 'role.name', to_jsonb('Administrator'::TEXT), FALSE)
    , ('instance1', 'developer', 'role.name', to_jsonb('Developer'::TEXT), FALSE)
    , ('instance1', 'user', 'role.name', to_jsonb('User'::TEXT), FALSE)

    , ('instance1', 'grant-org1', 'grantee', to_jsonb('org1'::TEXT), TRUE)
    , ('instance1', 'grant-org1', 'roles', to_jsonb(ARRAY['admin','developer']::TEXT[]), FALSE)
    , ('instance1', 'grant-org2', 'grantee', to_jsonb('org2'::TEXT), TRUE)
    , ('instance1', 'grant-org2', 'roles', to_jsonb(ARRAY['user']::TEXT[]), FALSE)
;