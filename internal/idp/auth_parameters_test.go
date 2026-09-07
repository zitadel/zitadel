package idp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveAuthorizationParameters(t *testing.T) {
	type args struct {
		static  map[string]string
		forward []string
		params  []Parameter
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "no configuration",
			args: args{
				params: []Parameter{AuthorizationParameters{"acr_values": "loa3"}},
			},
			want: map[string]string{},
		},
		{
			name: "static parameters only",
			args: args{
				static: map[string]string{"acr_values": "loa3", "prompt": "login"},
			},
			want: map[string]string{"acr_values": "loa3", "prompt": "login"},
		},
		{
			name: "only allow-listed parameters are forwarded",
			args: args{
				forward: []string{"acr_values", "max_age"},
				params: []Parameter{AuthorizationParameters{
					"acr_values": "loa3",
					"max_age":    "300",
					"ui_locales": "de",
				}},
			},
			want: map[string]string{"acr_values": "loa3", "max_age": "300"},
		},
		{
			name: "static parameters take precedence over forwarded ones",
			args: args{
				static:  map[string]string{"acr_values": "static"},
				forward: []string{"acr_values"},
				params:  []Parameter{AuthorizationParameters{"acr_values": "forwarded"}},
			},
			want: map[string]string{"acr_values": "static"},
		},
		{
			name: "reserved parameters are dropped",
			args: args{
				static:  map[string]string{"client_id": "evil", "State": "evil", "acr_values": "loa3"},
				forward: []string{"redirect_uri", "acr_values"},
				params: []Parameter{AuthorizationParameters{
					"redirect_uri": "https://evil.com",
				}},
			},
			want: map[string]string{"acr_values": "loa3"},
		},
		{
			name: "names are normalized",
			args: args{
				static:  map[string]string{"  ACR_Values ": "loa3"},
				forward: []string{" Max_Age "},
				params:  []Parameter{AuthorizationParameters{"MAX_AGE": "300"}},
			},
			want: map[string]string{"acr_values": "loa3", "max_age": "300"},
		},
		{
			name: "an empty value is kept to suppress a default",
			args: args{
				static: map[string]string{"prompt": ""},
			},
			want: map[string]string{"prompt": ""},
		},
		{
			name: "other parameter types are ignored",
			args: args{
				forward: []string{"acr_values"},
				params:  []Parameter{UserAgentID("agent"), LoginHintParam("user@example.com")},
			},
			want: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveAuthorizationParameters(tt.args.static, tt.args.forward, tt.args.params)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsReservedAuthorizationParameter(t *testing.T) {
	for _, key := range []string{"client_id", "REDIRECT_URI", " state ", "code_challenge"} {
		assert.True(t, IsReservedAuthorizationParameter(key), key)
	}
	for _, key := range []string{"acr_values", "prompt", "max_age", "ui_locales", "login_hint", ""} {
		assert.False(t, IsReservedAuthorizationParameter(key), key)
	}
}
