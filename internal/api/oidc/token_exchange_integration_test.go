//go:build integration

package oidc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/client/tokenexchange"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

func setTokenExchangeFeature(t *testing.T, value bool) {
	iamCTX := Tester.WithAuthorization(CTX, integration.IAMOwner)

	_, err := Tester.Client.FeatureV2.SetInstanceFeatures(iamCTX, &feature.SetInstanceFeaturesRequest{
		OidcTokenExchange: proto.Bool(value),
	})
	require.NoError(t, err)
	time.Sleep(time.Second)
}

func resetFeatures(t *testing.T) {
	iamCTX := Tester.WithAuthorization(CTX, integration.IAMOwner)
	_, err := Tester.Client.FeatureV2.ResetInstanceFeatures(iamCTX, &feature.ResetInstanceFeaturesRequest{})
	require.NoError(t, err)
	time.Sleep(time.Second)
}

func setImpersonationPolicy(t *testing.T, value bool) {
	iamCTX := Tester.WithAuthorization(CTX, integration.IAMOwner)

	policy, err := Tester.Client.Admin.GetSecurityPolicy(iamCTX, &admin.GetSecurityPolicyRequest{})
	require.NoError(t, err)
	if policy.GetPolicy().GetEnableImpersonation() != value {
		_, err = Tester.Client.Admin.SetSecurityPolicy(iamCTX, &admin.SetSecurityPolicyRequest{
			EnableImpersonation: value,
		})
		require.NoError(t, err)
	}
	time.Sleep(time.Second)
}

func createMachineUserPATWithMembership(t *testing.T, roles ...string) (userID, pat string) {
	iamCTX := Tester.WithAuthorization(CTX, integration.IAMOwner)
	userID, pat, err := Tester.CreateMachineUserPATWithMembership(iamCTX, roles...)
	require.NoError(t, err)
	return userID, pat
}

func accessTokenVerifier(ctx context.Context, server rs.ResourceServer, subject, actorSubject string) func(t *testing.T, token string) {
	return func(t *testing.T, token string) {
		resp, err := rs.Introspect[*oidc.IntrospectionResponse](ctx, server, token)
		require.NoError(t, err)
		assert.True(t, resp.Active)
		if subject != "" {
			assert.Equal(t, subject, resp.Subject)
		}
		if actorSubject != "" {
			require.NotNil(t, resp.Actor)
			assert.Equal(t, actorSubject, resp.Actor.Subject)
		}
	}
}

func idTokenVerifier(ctx context.Context, provider rp.RelyingParty, subject, actorSubject string) func(t *testing.T, token string) {
	return func(t *testing.T, token string) {
		verifier := provider.IDTokenVerifier()
		resp, err := rp.VerifyIDToken[*oidc.IDTokenClaims](ctx, token, verifier)
		require.NoError(t, err)
		if subject != "" {
			assert.Equal(t, subject, resp.Subject)
		}
		if actorSubject != "" {
			require.NotNil(t, resp.Actor)
			assert.Equal(t, actorSubject, resp.Actor.Subject)
		}
	}
}

func refreshTokenVerifier(ctx context.Context, provider rp.RelyingParty, subject, actorSubject string) func(t *testing.T, token string) {
	return func(t *testing.T, token string) {
		clientAssertion, err := client.SignedJWTProfileAssertion(provider.OAuthConfig().ClientID, []string{provider.Issuer()}, time.Hour, provider.Signer())
		require.NoError(t, err)
		tokens, err := rp.RefreshTokens[*oidc.IDTokenClaims](ctx, provider, token, clientAssertion, oidc.ClientAssertionTypeJWTAssertion)
		require.NoError(t, err)

		if subject != "" {
			assert.Equal(t, subject, tokens.IDTokenClaims.Subject)
		}
		if actorSubject != "" {
			require.NotNil(t, tokens.IDTokenClaims.Actor)
			assert.Equal(t, actorSubject, tokens.IDTokenClaims.Actor.Subject)
		}
		assert.NotEmpty(t, tokens.RefreshToken)
	}
}

func TestServer_TokenExchange(t *testing.T) {
	t.Cleanup(func() {
		resetFeatures(t)
		setImpersonationPolicy(t, false)
	})

	client, keyData, err := Tester.CreateOIDCTokenExchangeClient(CTX)
	require.NoError(t, err)
	signer, err := rp.SignerFromKeyFile(keyData)()
	require.NoError(t, err)
	exchanger, err := tokenexchange.NewTokenExchangerJWTProfile(CTX, Tester.OIDCIssuer(), client.GetClientId(), signer)
	require.NoError(t, err)

	time.Sleep(time.Second)

	iamUserID, iamImpersonatorPAT := createMachineUserPATWithMembership(t, "IAM_ADMIN_IMPERSONATOR")
	orgUserID, orgImpersonatorPAT := createMachineUserPATWithMembership(t, "ORG_ADMIN_IMPERSONATOR")
	serviceUserID, noPermPAT := createMachineUserPATWithMembership(t)

	// exchange some tokens for later use
	setTokenExchangeFeature(t, true)
	teResp, err := tokenexchange.ExchangeToken(CTX, exchanger, noPermPAT, oidc.AccessTokenType, "", "", nil, nil, nil, oidc.AccessTokenType)
	require.NoError(t, err)

	patScopes := oidc.SpaceDelimitedArray{"openid", "profile", "urn:zitadel:iam:user:metadata", "urn:zitadel:iam:user:resourceowner"}

	relyingParty, err := rp.NewRelyingPartyOIDC(CTX, Tester.OIDCIssuer(), client.GetClientId(), "", "", []string{"openid"}, rp.WithJWTProfile(rp.SignerFromKeyFile(keyData)))
	require.NoError(t, err)
	resourceServer, err := Tester.CreateResourceServerJWTProfile(CTX, keyData)
	require.NoError(t, err)

	type settings struct {
		tokenExchangeFeature bool
		impersonationPolicy  bool
	}
	type args struct {
		SubjectToken       string
		SubjectTokenType   oidc.TokenType
		ActorToken         string
		ActorTokenType     oidc.TokenType
		Resource           []string
		Audience           []string
		Scopes             []string
		RequestedTokenType oidc.TokenType
	}
	type result struct {
		issuedTokenType    oidc.TokenType
		tokenType          string
		expiresIn          uint64
		scopes             oidc.SpaceDelimitedArray
		verifyAccessToken  func(t *testing.T, token string)
		verifyRefreshToken func(t *testing.T, token string)
		verifyIDToken      func(t *testing.T, token string)
	}
	tests := []struct {
		name     string
		settings settings
		args     args
		want     result
		wantErr  bool
	}{
		{
			name: "feature disabled error",
			settings: settings{
				tokenExchangeFeature: false,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:     noPermPAT,
				SubjectTokenType: oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "unsupported resource parameter",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:     noPermPAT,
				SubjectTokenType: oidc.AccessTokenType,
				Resource:         []string{"https://example.com"},
			},
			wantErr: true,
		},
		{
			name: "invalid subject token",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:     "foo",
				SubjectTokenType: oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "EXCHANGE: access token to default",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:     noPermPAT,
				SubjectTokenType: oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, serviceUserID, ""),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, serviceUserID, ""),
			},
		},
		{
			name: "EXCHANGE: access token to access token",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, serviceUserID, ""),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, serviceUserID, ""),
			},
		},
		{
			name: "EXCHANGE: access token to JWT",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.JWTTokenType,
			},
			want: result{
				issuedTokenType:   oidc.JWTTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, serviceUserID, ""),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, serviceUserID, ""),
			},
		},
		{
			name: "EXCHANGE: access token to ID Token",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.IDTokenType,
			},
			want: result{
				issuedTokenType:   oidc.IDTokenType,
				tokenType:         "N_A",
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: idTokenVerifier(CTX, relyingParty, serviceUserID, ""),
				verifyIDToken: func(t *testing.T, token string) {
					assert.Empty(t, token)
				},
			},
		},
		{
			name: "EXCHANGE: refresh token not allowed",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       teResp.RefreshToken,
				SubjectTokenType:   oidc.RefreshTokenType,
				RequestedTokenType: oidc.IDTokenType,
			},
			wantErr: true,
		},
		{
			name: "EXCHANGE: alternate scope for refresh token",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				Scopes:             []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile"},
			},
			want: result{
				issuedTokenType:    oidc.AccessTokenType,
				tokenType:          oidc.BearerToken,
				expiresIn:          43100,
				scopes:             []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile"},
				verifyAccessToken:  accessTokenVerifier(CTX, resourceServer, serviceUserID, ""),
				verifyIDToken:      idTokenVerifier(CTX, relyingParty, serviceUserID, ""),
				verifyRefreshToken: refreshTokenVerifier(CTX, relyingParty, "", ""),
			},
		},
		{
			name: "EXCHANGE: access token, requested token type not supported error",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.RefreshTokenType,
			},
			wantErr: true,
		},
		{
			name: "EXCHANGE: access token, invalid audience",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				Audience:           []string{"foo", "bar"},
			},
			wantErr: true,
		},
		{
			name: "IMPERSONATION: subject: userID, actor: access token, policy disabled error",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  false,
			},
			args: args{
				SubjectToken:       User.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "IMPERSONATION: subject: userID, actor: access token, membership not found error",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken:       User.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         noPermPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "IAM IMPERSONATION: subject: userID, actor: access token, success",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken:       User.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         iamImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, User.GetUserId(), iamUserID),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, User.GetUserId(), iamUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: userID, actor: access token, success",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken:       User.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, User.GetUserId(), orgUserID),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, User.GetUserId(), orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: access token, actor: access token, success",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken:       teResp.AccessToken,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, serviceUserID, orgUserID),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, serviceUserID, orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: ID token, actor: access token, success",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken:       teResp.IDToken,
				SubjectTokenType:   oidc.IDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, serviceUserID, orgUserID),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, serviceUserID, orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: JWT, actor: access token, success",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken: func() string {
					token, err := crypto.Sign(&oidc.JWTTokenRequest{
						Issuer:    client.GetClientId(),
						Subject:   User.GetUserId(),
						Audience:  oidc.Audience{Tester.OIDCIssuer()},
						ExpiresAt: oidc.FromTime(time.Now().Add(time.Hour)),
						IssuedAt:  oidc.FromTime(time.Now().Add(-time.Second)),
					}, signer)
					require.NoError(t, err)
					return token
				}(),
				SubjectTokenType:   oidc.JWTTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(CTX, resourceServer, User.GetUserId(), orgUserID),
				verifyIDToken:     idTokenVerifier(CTX, relyingParty, User.GetUserId(), orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: access token, actor: access token, with refresh token, success",
			settings: settings{
				tokenExchangeFeature: true,
				impersonationPolicy:  true,
			},
			args: args{
				SubjectToken:       teResp.AccessToken,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
				Scopes:             []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
			},
			want: result{
				issuedTokenType:    oidc.AccessTokenType,
				tokenType:          oidc.BearerToken,
				expiresIn:          43100,
				scopes:             []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
				verifyAccessToken:  accessTokenVerifier(CTX, resourceServer, serviceUserID, orgUserID),
				verifyIDToken:      idTokenVerifier(CTX, relyingParty, serviceUserID, orgUserID),
				verifyRefreshToken: refreshTokenVerifier(CTX, relyingParty, serviceUserID, orgUserID),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTokenExchangeFeature(t, tt.settings.tokenExchangeFeature)
			setImpersonationPolicy(t, tt.settings.impersonationPolicy)

			got, err := tokenexchange.ExchangeToken(CTX, exchanger, tt.args.SubjectToken, tt.args.SubjectTokenType, tt.args.ActorToken, tt.args.ActorTokenType, tt.args.Resource, tt.args.Audience, tt.args.Scopes, tt.args.RequestedTokenType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.issuedTokenType, got.IssuedTokenType)
			assert.Equal(t, tt.want.tokenType, got.TokenType)
			assert.Greater(t, got.ExpiresIn, tt.want.expiresIn)
			assert.Equal(t, tt.want.scopes, got.Scopes)
			if tt.want.verifyAccessToken != nil {
				tt.want.verifyAccessToken(t, got.AccessToken)
			}
			if tt.want.verifyRefreshToken != nil {
				tt.want.verifyRefreshToken(t, got.RefreshToken)
			}
			if tt.want.verifyIDToken != nil {
				tt.want.verifyIDToken(t, got.IDToken)
			}
		})
	}
}

// This test tries to call the zitadel API with an impersonated token,
// which should fail.
func TestImpersonation_API_Call(t *testing.T) {
	client, keyData, err := Tester.CreateOIDCTokenExchangeClient(CTX)
	require.NoError(t, err)
	signer, err := rp.SignerFromKeyFile(keyData)()
	require.NoError(t, err)
	exchanger, err := tokenexchange.NewTokenExchangerJWTProfile(CTX, Tester.OIDCIssuer(), client.GetClientId(), signer)
	require.NoError(t, err)
	resourceServer, err := Tester.CreateResourceServerJWTProfile(CTX, keyData)
	require.NoError(t, err)

	setTokenExchangeFeature(t, true)
	setImpersonationPolicy(t, true)
	t.Cleanup(func() {
		resetFeatures(t)
		setImpersonationPolicy(t, false)
	})

	iamUserID, iamImpersonatorPAT := createMachineUserPATWithMembership(t, "IAM_ADMIN_IMPERSONATOR")
	iamOwner := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.IAMOwner)

	// impersonating the IAM owner!
	resp, err := tokenexchange.ExchangeToken(CTX, exchanger, iamOwner.Token, oidc.AccessTokenType, iamImpersonatorPAT, oidc.AccessTokenType, nil, nil, nil, oidc.AccessTokenType)
	require.NoError(t, err)
	accessTokenVerifier(CTX, resourceServer, iamOwner.ID, iamUserID)

	impersonatedCTX := Tester.WithAuthorizationToken(CTX, resp.AccessToken)
	_, err = Tester.Client.Admin.GetAllowedLanguages(impersonatedCTX, &admin.GetAllowedLanguagesRequest{})
	status := status.Convert(err)
	assert.Equal(t, codes.PermissionDenied, status.Code())
	assert.Equal(t, "Errors.TokenExchange.Token.NotForAPI (APP-Shi0J)", status.Message())
}
