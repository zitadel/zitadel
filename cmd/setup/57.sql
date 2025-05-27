DO $$ BEGIN
CREATE TYPE projections.object_level AS ENUM
    ('instance', 'organization');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS projections.resource_counts
(
    instance_id TEXT,
    table_name TEXT,
    owner_type projections.object_level DEFAULT 'instance',
    owner_id TEXT,
	resource_name TEXT NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    amount integer NOT NULL DEFAULT 1 CHECK (amount >= 0),
    
	PRIMARY KEY (instance_id, table_name, owner_type, owner_id)
);

CREATE OR REPLACE FUNCTION projections.parse_resource_count_pkey(
	tg_schema TEXT,
	tg_table_name TEXT,
	tg_argv TEXT[],
	tg_record RECORD,

	instance_id OUT TEXT,
	table_name OUT TEXT,
	owner_type OUT projections.object_level,
	owner_id OUT TEXT
)
	LANGUAGE 'plpgsql' IMMUTABLE
AS $$
DECLARE
	instance_id_column TEXT := COALESCE(tg_argv[2], 'instance_id');
	owner_id_column TEXT := COALESCE(tg_argv[3], 'resource_owner');
BEGIN
	table_name := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
	owner_type := COALESCE(tg_argv[1], 'instance');

	-- Get the ids from the RECORD
	EXECUTE format('SELECT ($1).%I, ($1).%I', instance_id_column, owner_id_column)
	INTO instance_id, owner_id
	USING tg_record;
END
$$;

-- count_resource is a trigger function which increases or decreases the count of a resource.
-- When creating the trigger the following arguments (TG_ARGV) can be passed:
-- 1. The name of the resource (required)
-- 1. The type of the owner, default is 'instance'
-- 2. The column name of the instance id, default is 'instance_id'
-- 3. The column name of the owner id, default is 'resource_owner'
CREATE OR REPLACE FUNCTION projections.count_resource()
    RETURNS trigger
    LANGUAGE 'plpgsql' VOLATILE
AS $$
BEGIN
	IF (TG_OP = 'INSERT') THEN
		INSERT INTO projections.resource_counts(instance_id, table_name, owner_type, owner_id, resource_name)
		SELECT projections.parse_resource_count_pkey(TG_TABLE_SCHEMA, TG_TABLE_NAME,TG_ARGV, NEW), TG_ARGV[0]	
		ON CONFLICT (instance_id, table_name, owner_type, owner_id) DO
		UPDATE SET updated_at = now(), amount = projections.resource_counts.amount + 1;
		RETURN NEW;
	ELSEIF (TG_OP = 'DELETE') THEN
		UPDATE projections.resource_counts
		SET updated_at = now(), amount = amount - 1
		WHERE (instance_id, table_name, owner_type, owner_id) = 
			projections.parse_resource_count_pkey(TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_ARGV, OLD)
		AND amount > 0; -- prevent check failure on negative amount.
		RETURN OLD;
	END IF;
END
$$;

-- delete_table_counts removes all resource counts for a TRUNCATEd table.
CREATE OR REPLACE FUNCTION projections.delete_table_counts()
	RETURNS trigger
	LANGUAGE 'plpgsql'
AS $$
DECLARE
	-- trigger variables
	tg_table_name TEXT := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
BEGIN
	DELETE FROM projections.resource_counts
	WHERE table_name = tg_table_name;
END
$$;

CREATE OR REPLACE FUNCTION projections.delete_resource_counts()
	RETURNS trigger
	LANGUAGE 'plpgsql'
AS $$
BEGIN
	DELETE FROM projections.resource_counts
	WHERE (instance_id, owner_type, owner_id) = (
		SELECT instance_id, owner_type, owner_id
		FROM projections.parse_resource_count_pkey(TG_TABLE_SCHEMA, TG_TABLE_NAME, TG_ARGV, OLD)
	);
END
$$;