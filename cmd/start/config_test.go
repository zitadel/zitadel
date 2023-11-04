package start

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestMustNewConfig(t *testing.T) {
	type args struct {
		yaml string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{{
		name: "features ok",
		args: args{yaml: `
DefaultInstance:
  Features:
  - FeatureLoginDefaultOrg: true
`},
		want: &Config{
			DefaultInstance: command.InstanceSetup{
				Features: map[domain.Feature]any{
					domain.FeatureLoginDefaultOrg: true,
				},
			},
		},
	}, {
		name: "membership types ok",
		args: args{yaml: `
SystemAPIUsers:
- superuser:
    Memberships:
    - MemberType: System
    - MemberType: Organization
    - MemberType: IAM
`},
		want: &Config{
			SystemAPIUsers: map[string]*authz.SystemAPIUser{
				"superuser": {
					Memberships: authz.Memberships{{
						MemberType: authz.MemberTypeSystem,
					}, {
						MemberType: authz.MemberTypeOrganization,
					}, {
						MemberType: authz.MemberTypeIAM,
					}},
				},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetConfigType("yaml")
			err := v.ReadConfig(strings.NewReader(`Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
` + tt.args.yaml))
			require.NoError(t, err)
			tt.want.Log = &logging.Config{Level: "info"}
			tt.want.Actions = &actions.Config{HTTP: actions.HTTPConfig{DenyList: []actions.AddressChecker{}}}
			require.NoError(t, tt.want.Log.SetLogger())
			got := MustNewConfig(v)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustNewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
