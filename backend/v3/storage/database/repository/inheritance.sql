CREATE TABLE objects (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE instances(
    name VARCHAR(50) NOT NULL
    , PRIMARY KEY (id)
) INHERITS (objects);

CREATE TRIGGER set_updated_at
BEFORE UPDATE
ON instances
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE instance_objects(
    instance_id INT NOT NULL
    , PRIMARY KEY (instance_id, id)
    -- as foreign keys are not inherited we need to define them on the child tables
    --, CONSTRAINT fk_instance FOREIGN KEY (instance_id) REFERENCES instances(id)
) INHERITS (objects);

CREATE TABLE orgs(
    name VARCHAR(50) NOT NULL
    , PRIMARY KEY (instance_id, id)
    , CONSTRAINT fk_instance FOREIGN KEY (instance_id) REFERENCES instances(id)
) INHERITS (instance_objects);

CREATE TRIGGER set_updated_at
BEFORE UPDATE
ON orgs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE org_objects(
    org_id INT NOT NULL
    , PRIMARY KEY (instance_id, org_id, id)
    -- as foreign keys are not inherited we need to define them on the child tables
    -- CONSTRAINT fk_org FOREIGN KEY (instance_id, org_id) REFERENCES orgs(instance_id, id),
    -- CONSTRAINT fk_instance FOREIGN KEY (instance_id) REFERENCES instances(id)
) INHERITS (instance_objects);

CREATE TABLE users (
    username VARCHAR(50) NOT NULL
    , PRIMARY KEY (instance_id, org_id, id)
    -- as foreign keys are not inherited we need to define them on the child tables
    -- , CONSTRAINT fk_org FOREIGN KEY (instance_id, org_id) REFERENCES orgs(instance_id, id)
    -- , CONSTRAINT fk_instances FOREIGN KEY (instance_id) REFERENCES instances(id)
) INHERITS (org_objects);

CREATE INDEX idx_users_username ON users(username);

CREATE TRIGGER set_updated_at
BEFORE UPDATE
ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE human_users(
    first_name VARCHAR(50)
    , last_name VARCHAR(50)
    , PRIMARY KEY (instance_id, org_id, id)
    -- CONSTRAINT fk_user FOREIGN KEY (instance_id, org_id, id) REFERENCES users(instance_id, org_id, id),
    , CONSTRAINT fk_org FOREIGN KEY (instance_id, org_id) REFERENCES orgs(instance_id, id)
    , CONSTRAINT fk_instances FOREIGN KEY (instance_id) REFERENCES instances(id)
) INHERITS (users);

CREATE INDEX idx_human_users_username ON human_users(username);

CREATE TRIGGER set_updated_at
BEFORE UPDATE
ON human_users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TABLE machine_users(
    description VARCHAR(50)
    , PRIMARY KEY (instance_id, org_id, id)
    -- , CONSTRAINT fk_user FOREIGN KEY (instance_id, org_id, id) REFERENCES users(instance_id, org_id, id)
    , CONSTRAINT fk_org FOREIGN KEY (instance_id, org_id) REFERENCES orgs(instance_id, id)
    , CONSTRAINT fk_instances FOREIGN KEY (instance_id) REFERENCES instances(id)
) INHERITS (users);

CREATE INDEX idx_machine_users_username ON machine_users(username);

CREATE TRIGGER set_updated_at
BEFORE UPDATE
ON machine_users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE VIEW users_view AS (
SELECT 
    id
    , created_at
    , updated_at
    , deleted_at
    , instance_id
    , org_id
    , username
    , tableoid::regclass::TEXT AS type
    , first_name
    , last_name
    , NULL AS description 
FROM
    human_users

UNION

SELECT 
    id
    , created_at
    , updated_at
    , deleted_at
    , instance_id
    , org_id
    , username
    , tableoid::regclass::TEXT AS type
    , NULL AS first_name
    , NULL AS last_name
    , description 
FROM
    machine_users
);