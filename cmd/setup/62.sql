ALTER TABLE projections.resource_counts
    DROP CONSTRAINT resource_counts_instance_id_parent_type_parent_id_table_nam_key;

ALTER TABLE projections.resource_counts
    ADD CONSTRAINT unique_resource
        UNIQUE (instance_id, parent_type, parent_id, table_name, resource_name);

-- count_resource is a trigger function which increases or decreases the count of a resource.
-- When creating the trigger the following required arguments (TG_ARGV) can be passed:
-- 1. The type of the parent
-- 2. The column name of the instance id
-- 3. The column name of the owner id
-- 4. The name of the resource
-- 5. (optional) 'UP' or 'DOWN' to indicate if an UPDATE should count up or down.
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
	tg_update_method TEXT := coalesce(TG_ARGV[4], '');

	tg_instance_id TEXT;
	tg_parent_id TEXT;

	select_ids TEXT := format('SELECT ($1).%I, ($1).%I', tg_instance_id_column, tg_parent_id_column);
BEGIN
	IF (TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' and tg_update_method = 'UP')) THEN
		EXECUTE select_ids INTO tg_instance_id, tg_parent_id USING NEW;
		
		INSERT INTO projections.resource_counts(instance_id, table_name, parent_type, parent_id, resource_name)
		VALUES (tg_instance_id, tg_table_name, tg_parent_type, tg_parent_id, tg_resource_name)	
		ON CONFLICT (instance_id, table_name, parent_type, parent_id, resource_name) DO
		UPDATE SET updated_at = now(), amount = projections.resource_counts.amount + 1;
		
		RETURN NEW;
	ELSEIF (TG_OP = 'DELETE' OR (TG_OP = 'UPDATE' and tg_update_method = 'DOWN')) THEN
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
