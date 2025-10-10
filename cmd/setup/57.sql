CREATE TABLE IF NOT EXISTS projections.resource_counts
(
    id SERIAL PRIMARY KEY, -- allows for easy pagination
    instance_id TEXT NOT NULL,
    table_name TEXT NOT NULL, -- needed for trigger matching, not in reports
    parent_type TEXT NOT NULL,
    parent_id TEXT NOT NULL,
    resource_name TEXT NOT NULL, -- friendly name for reporting
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    amount INTEGER NOT NULL DEFAULT 1 CHECK (amount >= 0),
    
	UNIQUE (instance_id, parent_type, parent_id, table_name)
);

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
	tg_parent_id_column TEXT := TG_ARGV[2];
	tg_resource_name TEXT := TG_ARGV[3];
	
	tg_instance_id TEXT;
	tg_parent_id TEXT;

	select_ids TEXT := format('SELECT ($1).%I, ($1).%I', tg_instance_id_column, tg_parent_id_column);
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

-- delete_table_counts removes all resource counts for a TRUNCATED table.
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

-- delete_parent_counts removes all resource counts for a deleted parent.
-- 1. The type of the parent
-- 2. The column name of the instance id
-- 3. The column name of the owner id
CREATE OR REPLACE FUNCTION projections.delete_parent_counts()
	RETURNS trigger
	LANGUAGE 'plpgsql'
AS $$
DECLARE
	-- trigger variables
	tg_parent_type TEXT := TG_ARGV[0];
	tg_instance_id_column TEXT := TG_ARGV[1];
	tg_parent_id_column TEXT := TG_ARGV[2];
	
	tg_instance_id TEXT;
	tg_parent_id TEXT;

	select_ids TEXT := format('SELECT ($1).%I, ($1).%I', tg_instance_id_column, tg_parent_id_column);
BEGIN
	EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING OLD;
	
	DELETE FROM projections.resource_counts
		WHERE instance_id = tg_instance_id
		AND parent_type = tg_parent_type
		AND parent_id = tg_parent_id;
	
	RETURN OLD;
END
$$;
