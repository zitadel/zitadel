package model

import (
	"reflect"
	"testing"
)

func TestGetRequiredGrantTypes(t *testing.T) {
	type args struct {
		oidcConfig OIDCConfig
	}
	tests := []struct {
		name   string
		args   args
		result []OIDCGrantType
	}{
		{
			name: "oidc response type code",
			args: args{
				oidcConfig: OIDCConfig{
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
				},
			},
			result: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
		},
		{
			name: "oidc response type id_token",
			args: args{
				oidcConfig: OIDCConfig{
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
				},
			},
			result: []OIDCGrantType{OIDCGrantTypeImplicit},
		},
		{
			name: "oidc response type id_token and id_token token",
			args: args{
				oidcConfig: OIDCConfig{
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken},
				},
			},
			result: []OIDCGrantType{OIDCGrantTypeImplicit},
		},
		{
			name: "oidc response type code, id_token and id_token token",
			args: args{
				oidcConfig: OIDCConfig{
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken},
				},
			},
			result: []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.oidcConfig.getRequiredGrantTypes()
			if !reflect.DeepEqual(tt.result, result) {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func TestContainsOIDCGrantType(t *testing.T) {
	type args struct {
		grantTypes []OIDCGrantType
		grantType  OIDCGrantType
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "contains grant type",
			args: args{
				grantTypes: []OIDCGrantType{
					OIDCGrantTypeAuthorizationCode,
					OIDCGrantTypeImplicit,
				},
				grantType: OIDCGrantTypeImplicit,
			},
			result: true,
		},
		{
			name: "doesnt contain grant type",
			args: args{
				grantTypes: []OIDCGrantType{
					OIDCGrantTypeAuthorizationCode,
					OIDCGrantTypeRefreshToken,
				},
				grantType: OIDCGrantTypeImplicit,
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsOIDCGrantType(tt.args.grantTypes, tt.args.grantType)
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func TestUrlsAreHttps(t *testing.T) {
	type args struct {
		uris []string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "only https uris",
			args: args{
				uris: []string{
					"https://zitadel.ch",
					"https://caos.ch",
				},
			},
			result: true,
		},
		{
			name: "http localhost uris",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://localhost",
				},
			},
			result: false,
		},
		{
			name: "http not localhsot",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://caos.ch",
				},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := urlsAreHttps(tt.args.uris)
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func TestOnlyLocalhostIsHttp(t *testing.T) {
	type args struct {
		uris []string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{

		{
			name: "http not localhost",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://caos.ch",
				},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := onlyLocalhostIsHttp(tt.args.uris)
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
