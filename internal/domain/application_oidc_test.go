package domain

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

func TestApplicationValid(t *testing.T) {
	type args struct {
		app *OIDCApp
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "invalid clock skew",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "AppName",
					ClockSkew:     time.Minute * 1,
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "invalid clock skew minus",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "AppName",
					ClockSkew:     time.Minute * -1,
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype code",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype code",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: responsetype token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				},
			},
			result: false,
		},
		{
			name: "valid oidc application: responsetype code & id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "valid oidc application: responsetype code & id_token & token_id_token",
			args: args{
				app: &OIDCApp{
					ObjectRoot:    models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:         "AppID",
					AppName:       "Name",
					ResponseTypes: []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken},
					GrantTypes:    []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
				},
			},
			result: true,
		},
		{
			name: "invalid oidc application: invalid origin",
			args: args{
				app: &OIDCApp{
					ObjectRoot:        models.ObjectRoot{AggregateID: "AggregateID"},
					AppID:             "AppID",
					AppName:           "Name",
					ResponseTypes:     []OIDCResponseType{OIDCResponseTypeCode, OIDCResponseTypeIDToken, OIDCResponseTypeIDTokenToken},
					GrantTypes:        []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeImplicit},
					AdditionalOrigins: []string{"https://test.com/test"},
				},
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.app.IsValid()
			if result != tt.result {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func TestGetOIDCV1Compliance(t *testing.T) {
	type args struct {
		appType      OIDCApplicationType
		grantTypes   []OIDCGrantType
		authMethod   OIDCAuthMethodType
		redirectUris []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "none compliant",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOIDCV1Compliance(tt.args.appType, tt.args.grantTypes, tt.args.authMethod, tt.args.redirectUris)
			if !got.NoneCompliant {
				t.Error("compliance should be none compliant")
			}
			if len(got.Problems) == 0 || got.Problems[0] != "Application.OIDC.V1.NotCompliant" {
				t.Errorf("first entry of problems should be \"Application.OIDC.V1.NotCompliant\" but got %v", got.Problems)
			}
		})
	}
}

func Test_checkGrantTypesCombination(t *testing.T) {
	tests := []struct {
		name       string
		want       *Compliance
		grantTypes []OIDCGrantType
	}{
		{
			name:       "implicit",
			want:       new(Compliance),
			grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit},
		},
		{
			name: "refresh token and implicit",
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.GrantType.Refresh.NoAuthCode"},
			},
			grantTypes: []OIDCGrantType{OIDCGrantTypeImplicit, OIDCGrantTypeRefreshToken},
		},
		{
			name:       "device code flow and refresh token doesnt require OIDCGrantTypeImplicit",
			want:       &Compliance{},
			grantTypes: []OIDCGrantType{OIDCGrantTypeDeviceCode, OIDCGrantTypeRefreshToken},
		},
		{
			name:       "refresh token and authorization code",
			want:       &Compliance{},
			grantTypes: []OIDCGrantType{OIDCGrantTypeAuthorizationCode, OIDCGrantTypeRefreshToken},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compliance := new(Compliance)

			checkGrantTypesCombination(compliance, tt.grantTypes)

			if tt.want.NoneCompliant != compliance.NoneCompliant {
				t.Errorf("NoneCompliant: expected: %v, got %v", tt.want.NoneCompliant, compliance.NoneCompliant)
			}
			if !reflect.DeepEqual(tt.want.Problems, compliance.Problems) {
				t.Errorf("Problems: expected: %v, got %v", tt.want.Problems, compliance.Problems)
			}
		})
	}
}

func Test_checkRedirectURIs(t *testing.T) {
	type args struct {
		grantTypes   []OIDCGrantType
		appType      OIDCApplicationType
		redirectUris []string
	}
	tests := []struct {
		name string
		want *Compliance
		args args
	}{
		{
			name: "no redirect uris",
			want: &Compliance{
				NoneCompliant: true,
				Problems: []string{
					"Application.OIDC.V1.NoRedirectUris",
				},
			},
			args: args{},
		},
		{
			name: "implicit and authorization code",
			want: &Compliance{
				NoneCompliant: false,
				Problems:      []string{"Application.OIDC.V1.NotAllCombinationsAreAllowed"},
			},
			args: args{
				redirectUris: []string{"http://redirect.to/me"},
				grantTypes:   []OIDCGrantType{OIDCGrantTypeImplicit, OIDCGrantTypeAuthorizationCode},
			},
		},
		{
			name: "only implicit",
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed"},
			},
			args: args{
				redirectUris: []string{"http://redirect.to/me"},
				grantTypes:   []OIDCGrantType{OIDCGrantTypeImplicit},
				appType:      OIDCApplicationTypeUserAgent,
			},
		},
		{
			name: "only authorization code",
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb"},
			},
			args: args{
				redirectUris: []string{"http://redirect.to/me"},
				grantTypes:   []OIDCGrantType{OIDCGrantTypeAuthorizationCode},
				appType:      OIDCApplicationTypeUserAgent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compliance := new(Compliance)

			checkRedirectURIs(compliance, tt.args.grantTypes, tt.args.appType, tt.args.redirectUris)

			if tt.want.NoneCompliant != compliance.NoneCompliant {
				t.Errorf("NoneCompliant: expected: %v, got %v", tt.want.NoneCompliant, compliance.NoneCompliant)
			}
			if !reflect.DeepEqual(tt.want.Problems, compliance.Problems) {
				t.Errorf("Problems: expected: %v, got %v", tt.want.Problems, compliance.Problems)
			}
		})
	}
}

func Test_CheckRedirectUrisImplicitAndCode(t *testing.T) {
	type args struct {
		appType      OIDCApplicationType
		redirectUris []string
	}
	tests := []struct {
		name string
		want *Compliance
		args args
	}{
		{
			name: "implicit and code https",
			want: &Compliance{
				NoneCompliant: false,
				Problems:      nil,
			},
			args: args{
				redirectUris: []string{"https://redirect.to/me"},
			},
		},
		// {
		// 	name: "custom protocol, not native",
		// 	want: &Compliance{
		// 		NoneCompliant: true,
		// 		Problems:      []string{"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed"},
		// 	},
		// 	args: args{
		// 		redirectUris: []string{"protocol://redirect.to/me"},
		// 		appType:      OIDCApplicationTypeWeb,
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compliance := new(Compliance)

			CheckRedirectUrisImplicitAndCode(compliance, tt.args.appType, tt.args.redirectUris)

			if tt.want.NoneCompliant != compliance.NoneCompliant {
				t.Errorf("NoneCompliant: expected: %v, got %v", tt.want.NoneCompliant, compliance.NoneCompliant)
			}
			if !reflect.DeepEqual(tt.want.Problems, compliance.Problems) {
				t.Errorf("Problems: expected: %v, got %v", tt.want.Problems, compliance.Problems)
			}
		})
	}
}

func TestCheckRedirectUrisImplicitAndCode(t *testing.T) {
	type args struct {
		appType      OIDCApplicationType
		redirectUris []string
	}
	tests := []struct {
		name string
		args args
		want *Compliance
	}{
		{
			name: "only https",
			args: args{},
			want: &Compliance{},
		},
		{
			name: "custom protocol not native app",
			args: args{
				appType:      OIDCApplicationTypeWeb,
				redirectUris: []string{"custom://nirvana.com"},
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed"},
			},
		},
		{
			name: "http localhost user agent app",
			args: args{
				appType:      OIDCApplicationTypeUserAgent,
				redirectUris: []string{"http://localhost:9009"},
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb"},
			},
		},
		{
			name: "http, not only localhost native app",
			args: args{
				appType:      OIDCApplicationTypeNative,
				redirectUris: []string{"http://nirvana.com", "http://localhost:9009"},
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Native.RedirectUris.MustBeHttpLocalhost"},
			},
		},
		{
			name: "not allowed combination",
			args: args{
				appType:      OIDCApplicationTypeNative,
				redirectUris: []string{"https://nirvana.com", "cutom://nirvana.com"},
			},
			want: &Compliance{
				NoneCompliant: false,
				Problems:      []string{"Application.OIDC.V1.NotAllCombinationsAreAllowed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(Compliance)
			CheckRedirectUrisImplicitAndCode(got, tt.args.appType, tt.args.redirectUris)

			if tt.want.NoneCompliant != got.NoneCompliant {
				t.Errorf("NoneCompliant: expected: %v, got %v", tt.want.NoneCompliant, got.NoneCompliant)
			}
			if !reflect.DeepEqual(tt.want.Problems, got.Problems) {
				t.Errorf("Problems: expected: %v, got %v", tt.want.Problems, got.Problems)
			}
		})
	}
}

func TestCheckRedirectUrisImplicit(t *testing.T) {
	type args struct {
		appType      OIDCApplicationType
		redirectUris []string
	}
	tests := []struct {
		name string
		args args
		want *Compliance
	}{
		{
			name: "only https",
			args: args{},
			want: &Compliance{},
		},
		{
			name: "custom protocol",
			args: args{
				redirectUris: []string{"custom://nirvana.com"},
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Implicit.RedirectUris.CustomNotAllowed"},
			},
		},
		{
			name: "only http protocol, app type native, not only localhost",
			args: args{
				redirectUris: []string{"http://nirvana.com"},
				appType:      OIDCApplicationTypeNative,
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Native.RedirectUris.MustBeHttpLocalhost"},
			},
		},
		{
			name: "only http protocol, app type native, only localhost",
			args: args{
				redirectUris: []string{"http://localhost:8080"},
				appType:      OIDCApplicationTypeNative,
			},
			want: &Compliance{
				NoneCompliant: false,
				Problems:      nil,
			},
		},
		{
			name: "only http protocol, app type web",
			args: args{
				redirectUris: []string{"http://nirvana.com"},
				appType:      OIDCApplicationTypeWeb,
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Implicit.RedirectUris.HttpNotAllowed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(Compliance)
			CheckRedirectUrisImplicit(got, tt.args.appType, tt.args.redirectUris)

			if tt.want.NoneCompliant != got.NoneCompliant {
				t.Errorf("NoneCompliant: expected: %v, got %v", tt.want.NoneCompliant, got.NoneCompliant)
			}
			if !reflect.DeepEqual(tt.want.Problems, got.Problems) {
				t.Errorf("Problems: expected: %v, got %v", tt.want.Problems, got.Problems)
			}
		})
	}
}

func TestCheckRedirectUrisCode(t *testing.T) {
	type args struct {
		appType      OIDCApplicationType
		redirectUris []string
	}
	tests := []struct {
		name string
		args args
		want *Compliance
	}{
		{
			name: "only https",
			args: args{},
			want: &Compliance{},
		},
		{
			name: "custom prefix, app type web",
			args: args{
				redirectUris: []string{"custom://nirvana.com"},
				appType:      OIDCApplicationTypeWeb,
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Code.RedirectUris.CustomOnlyForNative"},
			},
		},
		{
			name: "only http protocol, app type user agent",
			args: args{
				redirectUris: []string{"http://nirvana.com"},
				appType:      OIDCApplicationTypeUserAgent,
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Code.RedirectUris.HttpOnlyForWeb"},
			},
		},
		{
			name: "only http protocol, app type native, only localhost",
			args: args{
				redirectUris: []string{"http://localhost:8080", "http://nirvana.com:8080"},
				appType:      OIDCApplicationTypeNative,
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Native.RedirectUris.MustBeHttpLocalhost"},
			},
		},
		{
			name: "custom protocol, not native",
			args: args{
				redirectUris: []string{"custom://nirvana.com"},
				appType:      OIDCApplicationTypeWeb,
			},
			want: &Compliance{
				NoneCompliant: true,
				Problems:      []string{"Application.OIDC.V1.Code.RedirectUris.CustomOnlyForNative"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(Compliance)
			CheckRedirectUrisCode(got, tt.args.appType, tt.args.redirectUris)

			if tt.want.NoneCompliant != got.NoneCompliant {
				t.Errorf("NoneCompliant: expected: %v, got %v", tt.want.NoneCompliant, got.NoneCompliant)
			}
			if !reflect.DeepEqual(tt.want.Problems, got.Problems) {
				t.Errorf("Problems: expected: %v, got %v", tt.want.Problems, got.Problems)
			}
		})
	}
}

func TestOIDCOriginAllowList(t *testing.T) {
	type args struct {
		redirectUris      []string
		additionalOrigins []string
	}
	type want struct {
		allowed []string
		err     func(error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no uris, no origins",
			args: args{},
			want: want{
				allowed: []string{},
			},
		},
		{
			name: "redirects invalid schema",
			args: args{
				redirectUris: []string{"https:// localhost:8080"},
			},
			want: want{
				allowed: nil,
				err: func(e error) bool {
					return strings.HasPrefix(e.Error(), "invalid chavalid character")
				},
			},
		},
		{
			name: "redirects additional",
			args: args{
				redirectUris: []string{"https://localhost:8080"},
			},
			want: want{
				allowed: []string{"https://localhost:8080"},
			},
		},
		{
			name: "additional origin",
			args: args{
				additionalOrigins: []string{"https://localhost:8080"},
			},
			want: want{
				allowed: []string{"https://localhost:8080"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, err := OIDCOriginAllowList(tt.args.redirectUris, tt.args.additionalOrigins)

			if tt.want.err == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tt.want.err == nil && err == nil {
				//ok
			} else if tt.want.err(err) {
				t.Errorf("unexpected err got %v", err)
			}

			if !reflect.DeepEqual(allowed, tt.want.allowed) {
				t.Errorf("expected list: %v, got: %v", tt.want.allowed, allowed)
			}
		})
	}
}
