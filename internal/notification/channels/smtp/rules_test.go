package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileRules(t *testing.T) {
	tests := []struct {
		name    string
		configs []RuleConfig
		wantErr string
	}{
		{
			name:    "no rules",
			configs: nil,
		},
		{
			name: "valid rule",
			configs: []RuleConfig{{
				Match:   RuleMatch{Hosts: []string{"smtp.example.com"}},
				Headers: []RuleHeader{{Name: "X-Instance-ID", Value: "{{.InstanceID}}"}},
			}},
		},
		{
			name: "invalid header name",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X Invalid", Value: "value"}},
			}},
			wantErr: `smtp rule 0: header "X Invalid": invalid name`,
		},
		{
			name: "header name with line break",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header\r\nBcc", Value: "value"}},
			}},
			wantErr: "invalid name",
		},
		{
			name: "reserved header name, case insensitive",
			configs: []RuleConfig{{}, {
				Headers: []RuleHeader{{Name: "reply-to", Value: "value"}},
			}},
			wantErr: `smtp rule 1: header "reply-to": reserved name`,
		},
		{
			name: "header value with line break",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: "value\r\nBcc: attacker@example.com"}},
			}},
			wantErr: "value must not contain line breaks",
		},
		{
			name: "invalid template",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: "{{.InstanceID"}},
			}},
			wantErr: `header "X-Header"`,
		},
		{
			name: "unsupported template, unknown field",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{.UserID}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, unknown field hidden by condition",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{if .InstanceID}}{{.UserID}}{{end}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, condition",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{if .InstanceID}}a{{else}}b{{end}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, function",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{printf "%s" .InstanceID}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, index function",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{index .InstanceID 100}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, pipeline",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{.InstanceID | len}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, chained field",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{.InstanceID.Foo}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, with",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{with .InstanceID}}{{.}}{{end}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, range",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{range .InstanceID}}{{.}}{{end}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, variable",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{$id := .InstanceID}}{{$id}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, dot",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{.}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, nested template",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{define "x"}}{{.UserID}}{{end}}{{template "x" .}}`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "unsupported template, defined template",
			configs: []RuleConfig{{
				Headers: []RuleHeader{{Name: "X-Header", Value: `{{define "x"}}{{.UserID}}{{end}}text`}},
			}},
			wantErr: `header "X-Header": value must only contain text and the placeholders {{.InstanceID}} and {{.OrgID}}`,
		},
		{
			name: "text and placeholders with trim markers",
			configs: []RuleConfig{{
				Headers: []RuleHeader{
					{Name: "X-Header", Value: "prefix-{{ .InstanceID }}-{{- .OrgID -}} -suffix"},
					{Name: "X-Empty", Value: ""},
					{Name: "X-Text", Value: "text"},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, err := CompileRules(tt.configs)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, rules, len(tt.configs))
		})
	}
}

func TestRules_Match(t *testing.T) {
	data := RuleData{InstanceID: "instance1", OrgID: "org1"}
	provider := &Config{
		SMTP: SMTP{Host: "smtp.example.com:587"},
		From: "noreply@example.com",
	}

	tests := []struct {
		name    string
		configs []RuleConfig
		config  *Config
		want    Rule
	}{
		{
			name:   "no rules",
			config: provider,
			want:   Rule{},
		},
		{
			name:    "no provider",
			configs: []RuleConfig{{DisableCustomHTML: true}},
			config:  nil,
			want:    Rule{},
		},
		{
			name:    "empty match applies to all providers",
			configs: []RuleConfig{{DisableCustomHTML: true}},
			config:  provider,
			want:    Rule{DisableCustomHTML: true},
		},
		{
			name: "host without port matches any port",
			configs: []RuleConfig{{
				Match:             RuleMatch{Hosts: []string{"smtp.example.com"}},
				DisableCustomHTML: true,
			}},
			config: provider,
			want:   Rule{DisableCustomHTML: true},
		},
		{
			name: "host with port matches",
			configs: []RuleConfig{{
				Match:             RuleMatch{Hosts: []string{"smtp.example.com:587"}},
				DisableCustomHTML: true,
			}},
			config: provider,
			want:   Rule{DisableCustomHTML: true},
		},
		{
			name: "host with other port does not match",
			configs: []RuleConfig{{
				Match:             RuleMatch{Hosts: []string{"smtp.example.com:25"}},
				DisableCustomHTML: true,
			}},
			config: provider,
			want:   Rule{},
		},
		{
			name: "host is case insensitive",
			configs: []RuleConfig{{
				Match:             RuleMatch{Hosts: []string{" SMTP.Example.com "}},
				DisableCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.EXAMPLE.com:587"}, From: "noreply@example.com"},
			want:   Rule{DisableCustomHTML: true},
		},
		{
			name: "provider host without port",
			configs: []RuleConfig{{
				Match:             RuleMatch{Hosts: []string{"smtp.example.com"}},
				DisableCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.example.com"}, From: "noreply@example.com"},
			want:   Rule{DisableCustomHTML: true},
		},
		{
			name: "sender domain matches",
			configs: []RuleConfig{{
				Match:             RuleMatch{SenderDomains: []string{"EXAMPLE.com"}},
				DisableCustomHTML: true,
			}},
			config: provider,
			want:   Rule{DisableCustomHTML: true},
		},
		{
			name: "subdomain of sender domain does not match",
			configs: []RuleConfig{{
				Match:             RuleMatch{SenderDomains: []string{"example.com"}},
				DisableCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.example.com:587"}, From: "noreply@mail.example.com"},
			want:   Rule{},
		},
		{
			name: "sender without domain does not match",
			configs: []RuleConfig{{
				Match:             RuleMatch{SenderDomains: []string{"example.com"}},
				DisableCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.example.com:587"}, From: "example.com"},
			want:   Rule{},
		},
		{
			name: "same host, own sender domain does not match",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					SenderDomains: []string{"example.com"},
				},
				DisableCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.example.com:587"}, From: "noreply@other.example"},
			want:   Rule{},
		},
		{
			name: "same sender domain, other host does not match",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					SenderDomains: []string{"example.com"},
				},
				DisableCustomHTML: true,
			}},
			config: &Config{SMTP: SMTP{Host: "smtp.other.example:587"}, From: "noreply@example.com"},
			want:   Rule{},
		},
		{
			name: "first match wins",
			configs: []RuleConfig{
				{
					Match:   RuleMatch{Hosts: []string{"smtp.other.example"}},
					Headers: []RuleHeader{{Name: "X-Rule", Value: "first"}},
				},
				{
					Match:   RuleMatch{Hosts: []string{"smtp.example.com"}},
					Headers: []RuleHeader{{Name: "X-Rule", Value: "second"}},
				},
				{
					Headers: []RuleHeader{{Name: "X-Rule", Value: "third"}},
				},
			},
			config: provider,
			want:   Rule{Headers: map[string]string{"X-Rule": "second"}},
		},
		{
			name: "headers are rendered",
			configs: []RuleConfig{{
				Match: RuleMatch{
					Hosts:         []string{"smtp.example.com"},
					SenderDomains: []string{"example.com"},
				},
				DisableCustomHTML: true,
				Headers: []RuleHeader{
					{Name: "X-Instance-ID", Value: "{{.InstanceID}}"},
					{Name: "X-Org-ID", Value: "{{.OrgID}}"},
					{Name: "X-Tag", Value: "{{.InstanceID}}-{{.OrgID}}"},
					{Name: "X-Static", Value: "static"},
				},
			}},
			config: provider,
			want: Rule{
				DisableCustomHTML: true,
				Headers: map[string]string{
					"X-Instance-ID": "instance1",
					"X-Org-ID":      "org1",
					"X-Tag":         "instance1-org1",
					"X-Static":      "static",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, err := CompileRules(tt.configs)
			require.NoError(t, err)
			got, err := rules.Match(tt.config, data)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
