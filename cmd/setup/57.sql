CREATE TABLE IF NOT EXISTS projections.resource_counts
(
    instance_id TEXT,
    table_name TEXT,
    parent_type TEXT,
    parent_id TEXT,
	resource_name TEXT NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    amount integer NOT NULL DEFAULT 1 CHECK (amount >= 0),
    
	PRIMARY KEY (instance_id, table_name, parent_type, parent_id)
);

CREATE OR REPLACE FUNCTION projections.parse_resource_count_pkey(
	tg_schema TEXT,
	tg_table_name TEXT,
	tg_argv TEXT[],
	tg_record RECORD,

	instance_id OUT TEXT,
	table_name OUT TEXT,
	parent_type OUT projections.object_level,
	parent_id OUT TEXT
)
	LANGUAGE 'plpgsql' IMMUTABLE
AS $$
DECLARE
	instance_id_column TEXT := COALESCE(tg_argv[2], 'instance_id');
	owner_id_column TEXT := COALESCE(tg_argv[3], 'resource_owner');
BEGIN
	table_name := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
	parent_type := COALESCE(tg_argv[1], 'instance');

	-- Get the ids from the RECORD
	EXECUTE format('SELECT ($1).%I, ($1).%I', instance_id_column, owner_id_column)
	INTO instance_id, parent_id
	USING tg_record;
END
$$;

-- count_resource is a trigger function which increases or decreases the count of a resource.
-- When creating the trigger the following required arguments (TG_ARGV) can be passed:
-- 1. The type of the parent
-- 2. The column name of the instance id
-- 3. The column name of the owner id
-- 4. The name of the resource
CREATE OR REPLACE FUNCTION projections.count_resource()
    RETURNS trigger
    LANGUAGE 'plpgsql' VOLATILE
AS $$
DECLARE
	-- trigger variables
	tg_table_name TEXT := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
	tg_parent_type TEXT := TG_ARGV[0];
	tg_instance_id_column TEXT := TG_ARGV[1];
	tg_owner_id_column TEXT := TG_ARGV[2];
	tg_resource_name TEXT := TG_ARGV[3];
	
	tg_instance_id TEXT;
	tg_parent_id TEXT;

	select_ids TEXT := format('SELECT ($1).%I, ($1).%I', tg_instance_id_column, tg_owner_id_column)
BEGIN
	IF (TG_OP = 'INSERT') THEN
		EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING NEW;
		
		INSERT INTO projections.resource_counts(instance_id, table_name, parent_type, parent_id, resource_name)
		VALUES (tg_instance_id, tg_table_name, tg_parent_type, tg_parent_id, tg_resource_name)	
		ON CONFLICT (instance_id, table_name, parent_type, parent_id) DO
		UPDATE SET updated_at = now(), amount = projections.resource_counts.amount + 1;
		
		RETURN NEW;
	ELSEIF (TG_OP = 'DELETE') THEN
		EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING OLD;

		UPDATE projections.resource_counts
		SET updated_at = now(), amount = amount - 1
		WHERE instance_id = tg_instance_id
			AND table_name = tg_table_name
			AND parent_type = tg_parent_type
			AND parent_id = tg_parent_id
			AND resource_name = tg_resource_name 
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

-- delete_resource_counts removes all resource counts for a deleted parent.
-- 1. The type of the parent
-- 2. The column name of the instance id
-- 3. The column name of the owner id
-- 4. The name of the resource
CREATE OR REPLACE FUNCTION projections.delete_resource_counts()
	RETURNS trigger
	LANGUAGE 'plpgsql'
AS $$
DECLARE
	-- trigger variables
	tg_table_name TEXT := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;
	tg_parent_type TEXT := TG_ARGV[0];
	tg_instance_id_column TEXT := TG_ARGV[1];
	tg_owner_id_column TEXT := TG_ARGV[2];
	
	tg_instance_id TEXT;
	tg_parent_id TEXT;

	select_ids TEXT := format('SELECT ($1).%I, ($1).%I', tg_instance_id_column, tg_owner_id_column)
BEGIN
	EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING OLD;
	
	DELETE FROM projections.resource_counts
	WHERE instance_id = tg_instance_id
		AND table_name = tg_table_name
		AND parent_type = tg_parent_type
		AND parent_id = tg_parent_id;
	
	RETURN OLD;
END
$$;