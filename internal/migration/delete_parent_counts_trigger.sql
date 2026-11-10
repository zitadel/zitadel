{{ define "delete_parent_counts_trigger" -}}

CREATE OR REPLACE TRIGGER delete_parent_counts_trigger
    AFTER DELETE
    ON {{ .Table }}
    FOR EACH ROW
    EXECUTE FUNCTION projections.delete_parent_counts(
        '{{ .ParentType }}', 
        '{{ .InstanceIDColumn }}', 
        '{{ .ParentIDColumn }}'
    );

{{- end -}}
