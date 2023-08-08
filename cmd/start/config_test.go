package start

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestMustNewConfig(t *testing.T) {

	encodedKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF6aStGRlNKTDdmNXl3NEtUd3pnTQpQMzRlUEd5Y20vTStrVDBNN1Y0Q2d4NVYzRWFESXZUUUtUTGZCYUVCNDV6YjlMdGpJWHpEdzByWFJvUzJoTzZ0CmgrQ1lRQ3ozS0N2aDA5QzBJenhaaUIySVMzSC9hVCs1Qng5RUZZK3ZuQWtaamNjYnlHNVlOUnZtdE9sbnZJZUkKSDdxWjB0RXdrUGZGNUdFWk5QSlB0bXkzVUdWN2lvZmRWUVMxeFJqNzMrYU13NXJ2SDREOElkeWlBQzNWZWtJYgpwdDBWajBTVVgzRHdLdG9nMzM3QnpUaVBrM2FYUkYwc2JGaFFvcWRKUkk4TnFnWmpDd2pxOXlmSTV0eXhZc3duCitKR3pIR2RIdlczaWRPRGxtd0V0NUsycGFzaVJJV0syT0dmcSt3MEVjbHRRSGFidXFFUGdabG1oQ2tSZE5maXgKQndJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		t.Fatal(err)
	}

	type want = struct {
		got    func(*Config) interface{}
		expect interface{}
	}
	tests := []struct {
		name string
		args map[string]interface{}
		want want
	}{{
		name: "Actions.HTTP.DenyList slice of strings",
		args: map[string]interface{}{
			"actions": map[string]interface{}{
				"http": map[string]interface{}{
					"denyList": []string{"localhost", "127.0.0.1", "foobar"},
				},
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.Actions.HTTP.DenyList
			},
			expect: []actions.AddressChecker{&actions.DomainChecker{Domain: "localhost"}, &actions.IPChecker{IP: net.ParseIP("127.0.0.1")}, &actions.DomainChecker{Domain: "foobar"}},
		},
	}, {
		name: "Actions.HTTP.DenyList string",
		args: map[string]interface{}{
			"actions": map[string]interface{}{
				"http": map[string]interface{}{
					"denyList": "localhost,127.0.0.1,foobar",
				},
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.Actions.HTTP.DenyList
			},
			expect: []actions.AddressChecker{&actions.DomainChecker{Domain: "localhost"}, &actions.IPChecker{IP: net.ParseIP("127.0.0.1")}, &actions.DomainChecker{Domain: "foobar"}},
		},
	}, {
		name: "SystemAPIUsers slice of users",
		args: map[string]interface{}{
			"systemApiUsers": []interface{}{
				map[string]interface{}{
					"systemuser": map[string]interface{}{
						"path": "/path/to/superuser/key.pem",
					},
					"systemuser2": map[string]interface{}{
						"keyData": encodedKey,
					},
				},
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.SystemAPIUsers
			},
			expect: map[string]*authz.SystemAPIUser{
				"systemuser": {
					Path: "/path/to/superuser/key.pem",
				},
				"systemuser2": {
					KeyData: decodedKey,
				},
			},
		},
	}, {
		name: "SystemAPIUsers string",
		args: map[string]interface{}{
			"systemApiUsers": fmt.Sprintf(`{"systemuser": {"path": "/path/to/superuser/key.pem"}, "systemuser2": {"keyData": "%s"}}`, encodedKey),
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.SystemAPIUsers
			},
			expect: map[string]*authz.SystemAPIUser{
				"systemuser": {
					Path: "/path/to/superuser/key.pem",
				},
				"systemuser2": {
					KeyData: decodedKey,
				},
			},
		},
	}, {
		name: "Telemetry.Headers map of strings or map of string slices",
		args: map[string]interface{}{
			"telemetry": map[string]interface{}{
				"headers": map[string]interface{}{
					"single-value": "single-value",
					"multi-value":  []string{"multi-value1", "multi-value2"},
				},
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.Telemetry.Headers
			},
			expect: http.Header{
				"single-value": []string{"single-value"},
				"multi-value":  []string{"multi-value1", "multi-value2"},
			},
		},
	}, {
		name: "Telemetry.Headers string",
		args: map[string]interface{}{
			"telemetry": map[string]interface{}{
				"headers": `{"single-value": ["single-value"], "multi-value": ["multi-value1", "multi-value2"]}`,
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.Telemetry.Headers
			},
			expect: http.Header{
				"single-value": []string{"single-value"},
				"multi-value":  []string{"multi-value1", "multi-value2"},
			},
		},
	}, {
		name: "DefaultInstance.MessageTexts slice of custom message texts",
		args: map[string]interface{}{
			"defaultInstance": map[string]interface{}{
				"messageTexts": []map[string]interface{}{{
					"messageTextType": "InitCode",
					"title":           "foo",
				}, {
					"messageTextType": "PasswordReset",
					"greeting":        "bar",
				}},
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.DefaultInstance.MessageTexts
			},
			expect: []*domain.CustomMessageText{{
				MessageTextType: "InitCode",
				Title:           "foo",
			}, {
				MessageTextType: "PasswordReset",
				Greeting:        "bar",
			}},
		},
	}, {
		name: "DefaultInstance.MessageTexts string",
		args: map[string]interface{}{
			"defaultInstance": map[string]interface{}{
				"messageTexts": `[{
					"messageTextType": "InitCode",
					"title":           "foo"
				}, {
					"messageTextType": "PasswordReset",
					"greeting":           "bar"
				}]`,
			},
		},
		want: want{
			got: func(c *Config) interface{} {
				return c.DefaultInstance.MessageTexts
			},
			expect: []*domain.CustomMessageText{{
				MessageTextType: "InitCode",
				Title:           "foo",
			}, {
				MessageTextType: "PasswordReset",
				Greeting:        "bar",
			}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			if err := v.MergeConfigMap(map[string]interface{}{
				"log":     &logging.Config{},
				"actions": &actions.Config{},
			}); err != nil {
				t.Fatal(err)
			}
			if err := v.MergeConfigMap(tt.args); err != nil {
				t.Fatal(err)
			}
			if got := tt.want.got(MustNewConfig(v)); !reflect.DeepEqual(got, tt.want.expect) {
				t.Errorf("MustNewConfig() = %v, want %v", got, tt.want.expect)
			}
		})
	}
}
