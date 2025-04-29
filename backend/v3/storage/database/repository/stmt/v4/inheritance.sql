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

CREATE TRIGGER set_updated_at
BEFORE UPDATE
ON machine_users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


select u.*, hu.first_name, hu.last_name, mu.description from users u
left join human_users hu on u.instance_id = hu.instance_id and u.org_id = hu.org_id and u.id = hu.id
left join machine_users mu on u.instance_id = mu.instance_id and u.org_id = mu.org_id and u.id = mu.id
-- where
--     u.instance_id = 1
--     and u.org_id = 3
--     and u.id = 7
;

create view users_view as (
SELECT 
    id
    , created_at
    , updated_at
    , deleted_at
    , instance_id
    , org_id
    , username
    , first_name
    , last_name
    , description 
FROM (
(SELECT 
    id
    , created_at
    , updated_at
    , deleted_at
    , instance_id
    , org_id
    , username
    , first_name
    , last_name
    , NULL AS description 
FROM
    human_users)

UNION

(SELECT 
    id
    , created_at
    , updated_at
    , deleted_at
    , instance_id
    , org_id
    , username
    , NULL AS first_name
    , NULL AS last_name
    , description 
FROM
    machine_users)
));