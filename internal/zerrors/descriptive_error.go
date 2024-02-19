package zerrors

import (
	"errors"
)

type DescriptiveError struct {
	Parent error
	Vars   map[string]interface{}
}

// Describe accepts variables which can be used to add data to error message templates using errors.As() and GetVars()
// It only prints the parents error message.
// In most cases, it is simpler to just wrap an error with additional information than replacing variables in templates for all the translations.
// Example template and resolved message with using the vars:
// Template: Instance not found{{ if .host }} by origin {{ .origin }} and host {{ .host }}. {{ if ( not ( eq .host .external_domain ) ) }}You might try the domain {{ .external_domain }} or ask{{ else }}Ask{{ end }} your ZITADEL provider to check out https://zitadel.com/docs/self-hosting/manage/custom-domain {{ end }}
// Message: Instance not found by origin http://127.0.0.1.sslip.io:8080 and host 127.0.0.1.sslip.io:8080. You might try the domain localhost or ask your ZITADEL provider to check out https://zitadel.com/docs/self-hosting/manage/custom-domain
func Describe(err error, vars map[string]interface{}) error {
	if err == nil {
		return nil
	}
	return &DescriptiveError{Vars: vars, Parent: err}
}

func (err *DescriptiveError) Error() string {
	if err.Parent != nil {
		return err.Parent.Error()
	}
	return ""
}

func (err *DescriptiveError) Unwrap() error {
	return err.Parent
}

// GetVars returns a flat map of all vars added by WithVars
func (err *DescriptiveError) GetVars() map[string]interface{} {
	var vars map[string]interface{}
	p := new(DescriptiveError)
	if err.Parent != nil {
		if errors.As(err.Parent, &p) {
			vars = p.GetVars()
		}
	}
	if vars == nil && err.Vars != nil {
		vars = make(map[string]interface{})
	}
	for k, v := range err.Vars {
		vars[k] = v
	}
	return vars
}
