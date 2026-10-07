package settings

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func Test_emailProviderRestrictionsToPb(t *testing.T) {
	rules, err := smtp.CompileRules([]smtp.RuleConfig{
		{
			Match: smtp.RuleMatch{Hosts: []string{"smtp.example.com"}, Users: []string{"token"}},
			RuleOptions: smtp.RuleOptions{
				RestrictCustomHTML:               true,
				SuppressReservedRecipientDomains: true,
				Limit:                            &smtp.RuleLimit{Count: 100, Window: smtp.RuleDuration(24 * time.Hour)},
			},
		},
		{
			Match:       smtp.RuleMatch{Hosts: []string{"smtp.other.example"}},
			RuleOptions: smtp.RuleOptions{SuppressReservedRecipientDomains: true},
		},
		{
			// a rule without options, e.g. only adding headers, does not restrict the provider
			Match:   smtp.RuleMatch{Hosts: []string{"smtp.headers.example"}},
			Headers: map[string]string{"X-Instance-ID": "{{.InstanceID}}"},
		},
	})
	require.NoError(t, err)
	provider := func(host, user string) *query.SMTPConfig {
		return &query.SMTPConfig{SMTPConfig: &query.SMTP{Host: host, User: user, SenderAddress: "noreply@example.com"}}
	}

	tests := []struct {
		name   string
		rules  smtp.Rules
		config *query.SMTPConfig
		want   *settings.EmailProviderRestrictions
	}{
		{
			name:   "no rules",
			rules:  nil,
			config: provider("smtp.example.com:587", "token"),
			want:   nil,
		},
		{
			name:   "http provider is not restricted",
			rules:  rules,
			config: &query.SMTPConfig{HTTPConfig: &query.HTTP{Endpoint: "https://relay.example.com"}},
			want:   nil,
		},
		{
			name:   "no rule matches",
			rules:  rules,
			config: provider("smtp.example.com:587", "other"),
			want:   nil,
		},
		{
			name:   "matching rule without options",
			rules:  rules,
			config: provider("smtp.headers.example:587", "token"),
			want:   nil,
		},
		{
			name:   "all restrictions",
			rules:  rules,
			config: provider("smtp.example.com:587", "token"),
			want: &settings.EmailProviderRestrictions{
				CustomHtmlRestricted:               true,
				ReservedRecipientDomainsSuppressed: true,
				SendingLimit: &settings.EmailProviderSendingLimit{
					Count:  100,
					Window: durationpb.New(24 * time.Hour),
				},
			},
		},
		{
			name:   "restriction without limit",
			rules:  rules,
			config: provider("smtp.other.example:25", ""),
			want: &settings.EmailProviderRestrictions{
				ReservedRecipientDomainsSuppressed: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := emailProviderRestrictionsToPb(tt.config, tt.rules)
			assert.Equal(t, tt.want, got)
		})
	}
}
