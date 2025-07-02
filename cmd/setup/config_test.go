package setup

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestMustNewConfig(t *testing.T) {
	encodedKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF6aStGRlNKTDdmNXl3NEtUd3pnTQpQMzRlUEd5Y20vTStrVDBNN1Y0Q2d4NVYzRWFESXZUUUtUTGZCYUVCNDV6YjlMdGpJWHpEdzByWFJvUzJoTzZ0CmgrQ1lRQ3ozS0N2aDA5QzBJenhaaUIySVMzSC9hVCs1Qng5RUZZK3ZuQWtaamNjYnlHNVlOUnZtdE9sbnZJZUkKSDdxWjB0RXdrUGZGNUdFWk5QSlB0bXkzVUdWN2lvZmRWUVMxeFJqNzMrYU13NXJ2SDREOElkeWlBQzNWZWtJYgpwdDBWajBTVVgzRHdLdG9nMzM3QnpUaVBrM2FYUkYwc2JGaFFvcWRKUkk4TnFnWmpDd2pxOXlmSTV0eXhZc3duCitKR3pIR2RIdlczaWRPRGxtd0V0NUsycGFzaVJJV0syT0dmcSt3MEVjbHRRSGFidXFFUGdabG1oQ2tSZE5maXgKQndJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		yaml string
	}
	tests := []struct {
		name string
		args args
		want func(*testing.T, *Config)
	}{{
		name: "features ok",
		args: args{yaml: `
DefaultInstance:
  Features:
    LoginDefaultOrg: true
    UserSchema: true
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.DefaultInstance.Features, &command.InstanceFeatures{
				LoginDefaultOrg: gu.Ptr(true),
				UserSchema:      gu.Ptr(true),
			})
		},
	}, {
		name: "system api users ok",
		args: args{yaml: `
SystemAPIUsers:
- superuser:
    Memberships:
    - MemberType: System
    - MemberType: Organization
    - MemberType: IAM
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.SystemAPIUsers, map[string]*authz.SystemAPIUser{
				"superuser": {
					Memberships: authz.Memberships{{
						MemberType: authz.MemberTypeSystem,
					}, {
						MemberType: authz.MemberTypeOrganization,
					}, {
						MemberType: authz.MemberTypeIAM,
					}},
				},
			})
		},
	}, {
		name: "system api users string ok",
		args: args{yaml: fmt.Sprintf(`
SystemAPIUsers: >
  {"systemuser": {"path": "/path/to/superuser/key.pem"}, "systemuser2": {"keyData": "%s"}}
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`, encodedKey)},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.SystemAPIUsers, map[string]*authz.SystemAPIUser{
				"systemuser": {
					Path: "/path/to/superuser/key.pem",
				},
				"systemuser2": {
					KeyData: decodedKey,
				},
			})
		},
	}, {
		name: "headers ok",
		args: args{yaml: `
Telemetry:
  Headers:
    single-value: single-value
    multi-value:
    - multi-value1
    - multi-value2
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.Telemetry.Headers, http.Header{
				"single-value": []string{"single-value"},
				"multi-value":  []string{"multi-value1", "multi-value2"},
			})
		},
	}, {
		name: "headers string ok",
		args: args{yaml: `
Telemetry:
  Headers: >
    {"single-value": "single-value", "multi-value": ["multi-value1", "multi-value2"]}
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.Telemetry.Headers, http.Header{
				"single-value": []string{"single-value"},
				"multi-value":  []string{"multi-value1", "multi-value2"},
			})
		},
	}, {
		name: "message texts ok",
		args: args{yaml: `
DefaultInstance:
  MessageTexts:
  - MessageTextType: InitCode
    Title: foo
  - MessageTextType: PasswordReset
    Greeting: bar
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.DefaultInstance.MessageTexts, []*domain.CustomMessageText{{
				MessageTextType: "InitCode",
				Title:           "foo",
			}, {
				MessageTextType: "PasswordReset",
				Greeting:        "bar",
			}})
		},
	}, {
		name: "message texts string ok",
		args: args{yaml: `
DefaultInstance:
  MessageTexts: >
    [{"messageTextType": "InitCode", "title": "foo"}, {"messageTextType": "PasswordReset", "greeting": "bar"}]
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.DefaultInstance.MessageTexts, []*domain.CustomMessageText{{
				MessageTextType: "InitCode",
				Title:           "foo",
			}, {
				MessageTextType: "PasswordReset",
				Greeting:        "bar",
			}})
		},
	}, {
		name: "roles ok",
		args: args{yaml: `
InternalAuthZ:
  RolePermissionMappings:
  - Role: IAM_OWNER
    Permissions:
    - iam.write
  - Role: ORG_OWNER
    Permissions:
    - org.write
    - org.read
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.InternalAuthZ, authz.Config{
				RolePermissionMappings: []authz.RoleMapping{
					{Role: "IAM_OWNER", Permissions: []string{"iam.write"}},
					{Role: "ORG_OWNER", Permissions: []string{"org.write", "org.read"}},
				},
			})
		},
	}, {
		name: "roles string ok",
		args: args{yaml: `
InternalAuthZ:
  RolePermissionMappings: >
    [{"role": "IAM_OWNER", "permissions": ["iam.write"]}, {"role": "ORG_OWNER", "permissions": ["org.write", "org.read"]}]
Log:
  Level: info
Actions:
  HTTP:
    DenyList: []
`},
		want: func(t *testing.T, config *Config) {
			assert.Equal(t, config.InternalAuthZ, authz.Config{
				RolePermissionMappings: []authz.RoleMapping{
					{Role: "IAM_OWNER", Permissions: []string{"iam.write"}},
					{Role: "ORG_OWNER", Permissions: []string{"org.write", "org.read"}},
				},
			})
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			v.SetConfigType("yaml")
			require.NoError(t, v.ReadConfig(strings.NewReader(tt.args.yaml)))
			got := MustNewConfig(v)
			tt.want(t, got)
		})
	}
}
