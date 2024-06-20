package oidc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
)

func TestServer_createDiscoveryConfig(t *testing.T) {
	type fields struct {
		LegacyServer        *op.LegacyServer
		signingKeyAlgorithm string
	}
	type args struct {
		ctx                context.Context
		supportedUILocales []language.Tag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *oidc.DiscoveryConfiguration
	}{
		{
			"config",
			fields{
				LegacyServer: op.NewLegacyServer(
					func() *op.Provider {
						provider, _ := op.NewForwardedOpenIDProvider("path",
							&op.Config{
								CodeMethodS256:          true,
								AuthMethodPost:          true,
								AuthMethodPrivateKeyJWT: true,
								GrantTypeRefreshToken:   true,
								RequestObjectSupported:  true,
							},
							nil,
						)
						return provider
					}(),
					op.Endpoints{
						Authorization:       op.NewEndpoint("auth"),
						Token:               op.NewEndpoint("token"),
						Introspection:       op.NewEndpoint("introspect"),
						Userinfo:            op.NewEndpoint("userinfo"),
						Revocation:          op.NewEndpoint("revoke"),
						EndSession:          op.NewEndpoint("logout"),
						JwksURI:             op.NewEndpoint("keys"),
						DeviceAuthorization: op.NewEndpoint("device"),
					},
				),
				signingKeyAlgorithm: "RS256",
			},
			args{
				ctx:                op.ContextWithIssuer(context.Background(), "https://issuer.com"),
				supportedUILocales: []language.Tag{language.English, language.German},
			},
			&oidc.DiscoveryConfiguration{
				Issuer:                                             "https://issuer.com",
				AuthorizationEndpoint:                              "https://issuer.com/auth",
				TokenEndpoint:                                      "https://issuer.com/token",
				IntrospectionEndpoint:                              "https://issuer.com/introspect",
				UserinfoEndpoint:                                   "https://issuer.com/userinfo",
				RevocationEndpoint:                                 "https://issuer.com/revoke",
				EndSessionEndpoint:                                 "https://issuer.com/logout",
				DeviceAuthorizationEndpoint:                        "https://issuer.com/device",
				CheckSessionIframe:                                 "",
				JwksURI:                                            "https://issuer.com/keys",
				RegistrationEndpoint:                               "",
				ScopesSupported:                                    []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopePhone, oidc.ScopeAddress, oidc.ScopeOfflineAccess},
				ResponseTypesSupported:                             []string{string(oidc.ResponseTypeCode), string(oidc.ResponseTypeIDTokenOnly), string(oidc.ResponseTypeIDToken)},
				ResponseModesSupported:                             []string{string(oidc.ResponseModeQuery), string(oidc.ResponseModeFragment), string(oidc.ResponseModeFormPost)},
				GrantTypesSupported:                                []oidc.GrantType{oidc.GrantTypeCode, oidc.GrantTypeImplicit, oidc.GrantTypeRefreshToken, oidc.GrantTypeBearer},
				ACRValuesSupported:                                 nil,
				SubjectTypesSupported:                              []string{"public"},
				IDTokenSigningAlgValuesSupported:                   []string{"RS256"},
				IDTokenEncryptionAlgValuesSupported:                nil,
				IDTokenEncryptionEncValuesSupported:                nil,
				UserinfoSigningAlgValuesSupported:                  nil,
				UserinfoEncryptionAlgValuesSupported:               nil,
				UserinfoEncryptionEncValuesSupported:               nil,
				RequestObjectSigningAlgValuesSupported:             []string{"RS256"},
				RequestObjectEncryptionAlgValuesSupported:          nil,
				RequestObjectEncryptionEncValuesSupported:          nil,
				TokenEndpointAuthMethodsSupported:                  []oidc.AuthMethod{oidc.AuthMethodNone, oidc.AuthMethodBasic, oidc.AuthMethodPost, oidc.AuthMethodPrivateKeyJWT},
				TokenEndpointAuthSigningAlgValuesSupported:         []string{"RS256"},
				RevocationEndpointAuthMethodsSupported:             []oidc.AuthMethod{oidc.AuthMethodNone, oidc.AuthMethodBasic, oidc.AuthMethodPost, oidc.AuthMethodPrivateKeyJWT},
				RevocationEndpointAuthSigningAlgValuesSupported:    []string{"RS256"},
				IntrospectionEndpointAuthMethodsSupported:          []oidc.AuthMethod{oidc.AuthMethodBasic, oidc.AuthMethodPrivateKeyJWT},
				IntrospectionEndpointAuthSigningAlgValuesSupported: []string{"RS256"},
				DisplayValuesSupported:                             nil,
				ClaimTypesSupported:                                nil,
				ClaimsSupported:                                    []string{"sub", "aud", "exp", "iat", "iss", "auth_time", "nonce", "acr", "amr", "c_hash", "at_hash", "act", "scopes", "client_id", "azp", "preferred_username", "name", "family_name", "given_name", "locale", "email", "email_verified", "phone_number", "phone_number_verified"},
				ClaimsParameterSupported:                           false,
				CodeChallengeMethodsSupported:                      []oidc.CodeChallengeMethod{"S256"},
				ServiceDocumentation:                               "",
				ClaimsLocalesSupported:                             nil,
				UILocalesSupported:                                 []language.Tag{language.English, language.German},
				RequestParameterSupported:                          true,
				RequestURIParameterSupported:                       false,
				RequireRequestURIRegistration:                      false,
				OPPolicyURI:                                        "",
				OPTermsOfServiceURI:                                "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				LegacyServer:        tt.fields.LegacyServer,
				signingKeyAlgorithm: tt.fields.signingKeyAlgorithm,
			}
			assert.Equalf(t, tt.want, s.createDiscoveryConfig(tt.args.ctx, tt.args.supportedUILocales), "createDiscoveryConfig(%v)", tt.args.ctx)
		})
	}
}
