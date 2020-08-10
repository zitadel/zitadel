package model

import (
	"reflect"
	"testing"
)

func TestGetOIDCC1Compliance(t *testing.T) {
	type args struct {
		appType      OIDCApplicationType
		grantTypes   []OIDCGrantType
		authMethod   OIDCAuthMethodType
		redirectUris []string
	}
	type result struct {
		noneCompliant      bool
		complianceProblems []string
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			name: "Native: codeflow custom redirect (compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Native: codeflow http redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "Native: codeflow http://localhost redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "Native: codeflow http://localhost: redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost:1234/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "Native: codeflow https redirect (compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Native: codeflow invalid authmethod type (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypePost,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Native.AuthMethodType.NotNone",
				},
			},
		},
		{
			name: "Native: implicit custom redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed",
				},
			},
		},
		{
			name: "Native: implicit http redirect uri (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.NativeShouldBeHttpLocalhost",
				},
			},
		},
		{
			name: "Native: implicit http://localhost redirect uri (compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Native: implicit http://localhost: redirect uri (compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost:1234/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Native: implicit https redirect uri (compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Native: implicit and code (compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Native: implicit and code (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeNative,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "Web: code https redirect uri (compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Web: code http redirect uri (compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Web: code custom redirect uri (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.CustomOnlyForNative",
				},
			},
		},
		{
			name: "Web: implicit https uri (compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Web: implicit http redirect uri (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed",
				},
			},
		},
		{
			name: "Web: implicit custom redirect uri (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed",
				},
			},
		},
		{
			name: "Web: implicit http://localhost redirect uri (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed",
				},
			},
		},
		{
			name: "Web: implicit http://localhost: redirect uri (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost:1234/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed",
				},
			},
		},
		{
			name: "Web: implicit and code (compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit, OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "Web: implicit and code (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeWeb,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit, OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed",
				},
			},
		},
		{
			name: "UserAgent: code https redirect (compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "UserAgent: code http redirect (not compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "UserAgent: code http:localhost redirect (not compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "UserAgent: code http:localhost redirect (not compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost:1234/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb",
				},
			},
		},
		{
			name: "UserAgent: code custom redirect (not compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Code.RedirectUris.CustomOnlyForNative",
				},
			},
		},
		{
			name: "UserAgent: code authmethod type not none (not compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypePost,
				redirectUris: []string{
					"https://zitadel.chauth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.UserAgent.AuthMethodType.NotNone",
				},
			},
		},
		{
			name: "UserAgent: implicit https redirect (compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "UserAgent: implicit http redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed",
				},
			},
		},
		{
			name: "UserAgent: implicit custom redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed",
				},
			},
		},
		{
			name: "UserAgent: implicit http://localhost redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed",
				},
			},
		},
		{
			name: "UserAgent: implicit http://localhost: redirect (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"http://localhost:1234/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed",
				},
			},
		},
		{
			name: "UserAgent: implicit auth method not none (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
				authMethod: OIDCAuthMethodTypePost,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.UserAgent.AuthMethodType.NotNone",
				},
			},
		},
		{
			name: "UserAgent: implicit and code (compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit, OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
				},
			},
			result: result{
				noneCompliant: false,
			},
		},
		{
			name: "UserAgent: implicit and code (none compliant)",
			args: args{
				appType:    OIDCApplicationTypeUserAgent,
				grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit, OIDCGrantTypeAuthorizationCode},
				authMethod: OIDCAuthMethodTypeNone,
				redirectUris: []string{
					"https://zitadel.ch/auth/callback",
					"zitadel://auth/callback",
				},
			},
			result: result{
				noneCompliant: true,
				complianceProblems: []string{
					"Application.OIDC.V1.NotCompliant",
					"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOIDCV1Compliance(tt.args.appType, tt.args.grantTypes, tt.args.authMethod, tt.args.redirectUris)
			if tt.result.noneCompliant != result.NoneCompliant {
				t.Errorf("got wrong result nonecompliant: expected: %v, actual: %v ", tt.result.noneCompliant, result.NoneCompliant)
			}
			if tt.result.noneCompliant {
				if len(tt.result.complianceProblems) != len(result.Problems) {
					t.Errorf("got wrong result compliance problems len: expected: %v, actual: %v ", len(tt.result.complianceProblems), len(result.Problems))
				}
				if !reflect.DeepEqual(tt.result.complianceProblems, result.Problems) {
					t.Errorf("got wrong result compliance problems: expected: %v, actual: %v ", tt.result.complianceProblems, result.Problems)
				}
			}
		})
	}
}

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
		{
			name: "http localhost/",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://localhost/auth/callback",
				},
			},
			result: true,
		},
		{
			name: "http not localhost:",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://localhost:9090",
				},
			},
			result: true,
		},
		{
			name: "http not localhost",
			args: args{
				uris: []string{
					"https://zitadel.com",
					"http://localhost:9090",
					"http://zitadel.ch",
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
