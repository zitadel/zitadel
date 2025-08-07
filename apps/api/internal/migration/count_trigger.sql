{{ define "count_trigger" -}}
CREATE OR REPLACE TRIGGER count_{{ .Resource }}
    AFTER INSERT OR DELETE
    ON {{ .Table }}
    FOR EACH ROW
    EXECUTE FUNCTION projections.count_resource(
        '{{ .ParentType }}', 
        '{{ .InstanceIDColumn }}', 
        '{{ .ParentIDColumn }}',
        '{{ .Resource }}'
    );

CREATE OR REPLACE TRIGGER truncate_{{ .Resource }}_counts
    AFTER TRUNCATE
    ON {{ .Table }}
    FOR EACH STATEMENT
    EXECUTE FUNCTION projections.delete_table_counts();

-- Prevent inserts and deletes while we populate the counts.
LOCK TABLE {{ .Table }} IN SHARE MODE;

-- Populate the resource counts for the existing data in the table.
INSERT INTO projections.resource_counts(
	instance_id,
    table_name,
    parent_type,
    parent_id,
    resource_name,
    amount
)
SELECT
    {{ .InstanceIDColumn }},
    '{{ .Table }}',
    '{{ .ParentType }}',
    {{ .ParentIDColumn }},
    '{{ .Resource }}',
    COUNT(*) AS amount
FROM {{ .Table }}
GROUP BY ({{ .InstanceIDColumn }}, {{ .ParentIDColumn }})
ON CONFLICT (instance_id, table_name, parent_type, parent_id) DO
UPDATE SET updated_at = now(), amount = EXCLUDED.amount;

{{- end -}}
