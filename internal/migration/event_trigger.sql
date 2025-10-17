{{ define "event_trigger" -}}
-- In case the old triggers exist, drop them to avoid duplicates.
DROP TRIGGER IF EXISTS event_{{ .Resource }}_insert ON {{ .Table }};
DROP TRIGGER IF EXISTS event_{{ .Resource }}_update ON {{ .Table }};
DROP TRIGGER IF EXISTS event_{{ .Resource }}_delete ON {{ .Table }};

CREATE OR REPLACE TRIGGER event_{{ .Resource }}_insert
    AFTER INSERT
    ON {{ .Table }}
    FOR EACH ROW
    {{if .Conditions}}WHEN ({{.Conditions.ToSQL "NEW" true}}){{end}}
    EXECUTE FUNCTION analytics.record_service_ping_resource_event(
        '{{ .ParentType }}',
        '{{ .InstanceIDColumn }}',
        '{{ .ParentIDColumn }}'
    );

CREATE OR REPLACE TRIGGER event_{{ .Resource }}_update
    AFTER UPDATE
    ON {{ .Table }}
    FOR EACH ROW
    {{if .Conditions}}WHEN ({{.Conditions.ToSQL "NEW" true}} OR {{.Conditions.ToSQL "OLD" true}}){{end}}
    EXECUTE FUNCTION analytics.record_service_ping_resource_event(
        '{{ .ParentType }}',
        '{{ .InstanceIDColumn }}',
        '{{ .ParentIDColumn }}'
    );

CREATE OR REPLACE TRIGGER event_{{ .Resource }}_delete
    AFTER DELETE
    ON {{ .Table }}
    FOR EACH ROW
    {{if .Conditions}}WHEN ({{.Conditions.ToSQL "OLD" true}}){{end}}
    EXECUTE FUNCTION analytics.record_service_ping_resource_event(
        '{{ .ParentType }}',
        '{{ .InstanceIDColumn }}',
        '{{ .ParentIDColumn }}'
    );

{{- end -}}


