package mirror

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

func Test_newConfig_smtpRules(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{
			name: "yaml",
			yaml: `
Notifications:
  SMTPRules:
    - Match:
        Hosts:
          - smtp.example.com
      RestrictCustomHTML: true
      Headers:
        X-Instance-ID: "{{.InstanceID}}"
`,
		},
		{
			name: "string",
			yaml: `
Notifications:
  SMTPRules: >
    [{"Match": {"Hosts": ["smtp.example.com"]}, "RestrictCustomHTML": true, "Headers": {"X-Instance-ID": "{{.InstanceID}}"}}]
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetConfigType("yaml")
			require.NoError(t, v.ReadConfig(strings.NewReader(tt.yaml)))
			config := new(ProjectionsConfig)
			require.NoError(t, newConfig(v, config))

			// the rules are checked through the compiled result,
			// because keys of maps are lower-cased by viper when read from YAML, but not from the JSON string
			rules, err := smtp.CompileRules(config.Notifications.SMTPRules)
			require.NoError(t, err)
			provider := &smtp.Config{SMTP: smtp.SMTP{Host: "smtp.example.com:587"}}
			assert.Equal(t, smtp.Rule{
				RuleOptions: smtp.RuleOptions{
					RestrictCustomHTML: true,
				},
				Headers: map[string]string{"X-Instance-Id": "instance1"},
			}, rules.Match(provider, smtp.RuleData{InstanceID: "instance1"}))
		})
	}
}
