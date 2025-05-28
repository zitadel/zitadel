{{ template "delete_parent_counts" -}}

CREATE OR REPLACE TRIGGER delete_parent_counts
    AFTER DELETE
    ON {{ .Table }}
    FOR EACH ROW
    EXECUTE FUNCTION projections.delete_parent_counts(
        '{{ .ParentType }}', 
        '{{ .InstanceIDColumn }}', 
        '{{ .ParentIDColumn }}'
    );

{{- end -}}
