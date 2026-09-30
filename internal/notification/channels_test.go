package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

func Test_channels_SMTPRule(t *testing.T) {
	rules, err := smtp.CompileRules([]smtp.RuleConfig{{
		Match: smtp.RuleMatch{
			Hosts:         []string{"smtp.example.com"},
			SenderDomains: []string{"m.example.com"},
		},
		DisableCustomHTML: true,
		Headers: []smtp.RuleHeader{
			{Name: "X-Instance-ID", Value: "{{.InstanceID}}"},
			{Name: "X-Org-ID", Value: "{{.OrgID}}"},
		},
	}})
	require.NoError(t, err)

	tests := []struct {
		name   string
		rules  smtp.Rules
		config *smtp.Config
		want   smtp.Rule
	}{
		{
			name:   "no rules configured",
			rules:  nil,
			config: &smtp.Config{SMTP: smtp.SMTP{Host: "smtp.example.com:587"}, From: "noreply@m.example.com"},
			want:   smtp.Rule{},
		},
		{
			name:   "provider does not match",
			rules:  rules,
			config: &smtp.Config{SMTP: smtp.SMTP{Host: "smtp.example.com:587"}, From: "noreply@other.example"},
			want:   smtp.Rule{},
		},
		{
			name:   "provider matches, headers contain instance and org",
			rules:  rules,
			config: &smtp.Config{SMTP: smtp.SMTP{Host: "smtp.example.com:587"}, From: "noreply@m.example.com"},
			want: smtp.Rule{
				DisableCustomHTML: true,
				Headers: map[string]string{
					"X-Instance-ID": "instance1",
					"X-Org-ID":      "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &channels{smtpRules: tt.rules}
			ctx := authz.WithInstanceID(t.Context(), "instance1")
			got, err := c.SMTPRule(ctx, tt.config, "org1")
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
