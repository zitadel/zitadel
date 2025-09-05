{{ define "count_trigger" -}}
-- In case the old trigger exists, drop it to prevent duplicated counts.
DROP TRIGGER IF EXISTS count_{{ .Resource }} ON {{ .Table }};

CREATE OR REPLACE TRIGGER count_{{ .Resource }}_insert
    AFTER INSERT
    ON {{ .Table }}
    FOR EACH ROW
    -- Only count if the conditions are met in the newly added row.
    {{if .Conditions}}WHEN ({{.Conditions.ToSQL "NEW" true}}){{end}}
    EXECUTE FUNCTION projections.count_resource(
        '{{ .ParentType }}', 
        '{{ .InstanceIDColumn }}', 
        '{{ .ParentIDColumn }}',
        '{{ .Resource }}'
    );

CREATE OR REPLACE TRIGGER count_{{ .Resource }}_delete
    AFTER DELETE
    ON {{ .Table }}
    FOR EACH ROW
    -- Only count down if the conditions were met in the old / deleted row.
    {{if .Conditions}}WHEN ({{.Conditions.ToSQL "OLD" true}}){{end}}
    EXECUTE FUNCTION projections.count_resource(
        '{{ .ParentType }}',
        '{{ .InstanceIDColumn }}',
        '{{ .ParentIDColumn }}',
        '{{ .Resource }}'
    );

{{if .TrackChange}}
CREATE OR REPLACE TRIGGER count_{{ .Resource }}_update_up
    AFTER UPDATE
    ON {{ .Table }}
    FOR EACH ROW
    -- Only count up if the conditions are met in the new state, but were not in the old.
    WHEN ({{.Conditions.ToSQL "NEW" true}} AND {{.Conditions.ToSQL "OLD" false}})
    EXECUTE FUNCTION projections.count_resource(
        '{{ .ParentType }}',
        '{{ .InstanceIDColumn }}',
        '{{ .ParentIDColumn }}',
        '{{ .Resource }}',
        'UP'
    );

CREATE OR REPLACE TRIGGER count_{{ .Resource }}_update_down
    AFTER UPDATE
    ON {{ .Table }}
    FOR EACH ROW
    -- Only count down if the conditions are not met in the new state, but were in the old.
    WHEN ({{.Conditions.ToSQL "NEW" false}} AND {{.Conditions.ToSQL "OLD" true}})
    EXECUTE FUNCTION projections.count_resource(
        '{{ .ParentType }}',
        '{{ .InstanceIDColumn }}',
        '{{ .ParentIDColumn }}',
        '{{ .Resource }}',
        'DOWN'
    );
{{end}}

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
{{if .Conditions}}WHERE {{.Conditions.ToSQL .Table true}}{{end}}
GROUP BY ({{ .InstanceIDColumn }}, {{ .ParentIDColumn }})
ON CONFLICT (instance_id, table_name, parent_type, parent_id, resource_name) DO
UPDATE SET updated_at = now(), amount = EXCLUDED.amount;

{{- end -}}
