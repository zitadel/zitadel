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

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func setImpersonationPolicy(t *testing.T, instance *integration.Instance, value bool) {
	instance.SetImpersonationPolicy(CTX, t, value)
}

// awaitImpersonationDenied retries the token exchange until the impersonation is rejected
// because the actor lacks the required permission.
// Retries are needed because both the actor's permissions and the subject's admin state are
// resolved from eventually consistent projections. As long as the memberships are not projected
// yet, the exchange either succeeds or fails with an unrelated "membership not found" error.
func awaitImpersonationDenied(ctx context.Context, t *testing.T, exchanger tokenexchange.TokenExchanger, subjectUserID, actorToken string, requestedTokenType oidc.TokenType) {
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		_, err := tokenexchange.ExchangeToken(ctx, exchanger, subjectUserID, oidc_api.UserIDTokenType, actorToken, oidc.AccessTokenType, nil, nil, nil, requestedTokenType)
		if assert.Error(ttt, err) {
			assert.ErrorContains(ttt, err, "No matching permissions found")
		}
	},
		retryDuration,
		tick,
		"timed out waiting for the impersonation to be denied")
}

// awaitImpersonationAllowed retries the token exchange until it succeeds and returns the response.
// See [awaitImpersonationDenied] on why retries are needed.
func awaitImpersonationAllowed(ctx context.Context, t *testing.T, exchanger tokenexchange.TokenExchanger, subjectUserID, actorToken string, requestedTokenType oidc.TokenType) *oidc.TokenExchangeResponse {
	var resp *oidc.TokenExchangeResponse
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := tokenexchange.ExchangeToken(ctx, exchanger, subjectUserID, oidc_api.UserIDTokenType, actorToken, oidc.AccessTokenType, nil, nil, nil, requestedTokenType)
		if !assert.NoError(ttt, err) {
			return
		}
		resp = got
	},
		retryDuration,
		tick,
		"timed out waiting for the impersonation to be allowed")
	return resp
}

func createMachineUserPATWithMembership(ctx context.Context, t *testing.T, instance *integration.Instance, roles ...string) (userID, pat string) {
	userID, pat, err := instance.CreateMachineUserPATWithMembership(ctx, roles...)
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
	instance := integration.NewInstance(CTX)
	ctx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	userResp := instance.CreateHumanUser(ctx)

	client, keyData, err := instance.CreateOIDCTokenExchangeClient(ctx, t)
	require.NoError(t, err)
	signer, err := rp.SignerFromKeyFile(keyData)()
	require.NoError(t, err)
	exchanger, err := tokenexchange.NewTokenExchangerJWTProfile(ctx, instance.OIDCIssuer(), client.GetClientId(), signer)
	require.NoError(t, err)

	_, orgImpersonatorPAT := createMachineUserPATWithMembership(ctx, t, instance, "ORG_ADMIN_IMPERSONATOR")
	serviceUserID, noPermPAT := createMachineUserPATWithMembership(ctx, t, instance)

	teResp, err := tokenexchange.ExchangeToken(ctx, exchanger, noPermPAT, oidc.AccessTokenType, "", "", nil, nil, nil, oidc.AccessTokenType)
	require.NoError(t, err)

	patScopes := oidc.SpaceDelimitedArray{"openid", "profile", "urn:zitadel:iam:user:metadata", "urn:zitadel:iam:user:resourceowner"}

	relyingParty, err := rp.NewRelyingPartyOIDC(ctx, instance.OIDCIssuer(), client.GetClientId(), "", "", []string{"openid"}, rp.WithJWTProfile(rp.SignerFromKeyFile(keyData)))
	require.NoError(t, err)
	resourceServer, err := instance.CreateResourceServerJWTProfile(ctx, keyData)
	require.NoError(t, err)

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
		name    string
		args    args
		want    result
		wantErr bool
	}{
		{
			name: "unsupported resource parameter",
			args: args{
				SubjectToken:     noPermPAT,
				SubjectTokenType: oidc.AccessTokenType,
				Resource:         []string{"https://example.com"},
			},
			wantErr: true,
		},
		{
			name: "invalid subject token",
			args: args{
				SubjectToken:     "foo",
				SubjectTokenType: oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "EXCHANGE: access token to default",
			args: args{
				SubjectToken:     noPermPAT,
				SubjectTokenType: oidc.AccessTokenType,
			},
			want: result{
				issuedTokenType:   oidc.AccessTokenType,
				tokenType:         oidc.BearerToken,
				expiresIn:         43100,
				scopes:            patScopes,
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, serviceUserID, ""),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, serviceUserID, ""),
			},
		},
		{
			name: "EXCHANGE: access token to access token",
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, serviceUserID, ""),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, serviceUserID, ""),
			},
		},
		{
			name: "EXCHANGE: access token to JWT",
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, serviceUserID, ""),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, serviceUserID, ""),
			},
		},
		{
			name: "EXCHANGE: access token to ID Token",
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
				verifyAccessToken: idTokenVerifier(ctx, relyingParty, serviceUserID, ""),
				verifyIDToken: func(t *testing.T, token string) {
					assert.Empty(t, token)
				},
			},
		},
		{
			name: "EXCHANGE: refresh token not allowed",
			args: args{
				SubjectToken:       teResp.RefreshToken,
				SubjectTokenType:   oidc.RefreshTokenType,
				RequestedTokenType: oidc.IDTokenType,
			},
			wantErr: true,
		},
		{
			name: "EXCHANGE: alternate scope for refresh token",
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
				verifyAccessToken:  accessTokenVerifier(ctx, resourceServer, serviceUserID, ""),
				verifyIDToken:      idTokenVerifier(ctx, relyingParty, serviceUserID, ""),
				verifyRefreshToken: refreshTokenVerifier(ctx, relyingParty, "", ""),
			},
		},
		{
			name: "EXCHANGE: access token, requested token type not supported error",
			args: args{
				SubjectToken:       noPermPAT,
				SubjectTokenType:   oidc.AccessTokenType,
				RequestedTokenType: oidc.RefreshTokenType,
			},
			wantErr: true,
		},
		{
			name: "EXCHANGE: access token, invalid audience",
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
			args: args{
				SubjectToken:       userResp.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         orgImpersonatorPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenexchange.ExchangeToken(ctx, exchanger, tt.args.SubjectToken, tt.args.SubjectTokenType, tt.args.ActorToken, tt.args.ActorTokenType, tt.args.Resource, tt.args.Audience, tt.args.Scopes, tt.args.RequestedTokenType)
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

func TestServer_TokenExchangeImpersonation(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ctx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	userResp := instance.CreateHumanUser(ctx)

	// exchange some tokens for later use
	setImpersonationPolicy(t, instance, true)

	client, keyData, err := instance.CreateOIDCTokenExchangeClient(ctx, t)
	require.NoError(t, err)
	signer, err := rp.SignerFromKeyFile(keyData)()
	require.NoError(t, err)
	exchanger, err := tokenexchange.NewTokenExchangerJWTProfile(ctx, instance.OIDCIssuer(), client.GetClientId(), signer)
	require.NoError(t, err)

	iamUserID, iamImpersonatorPAT := createMachineUserPATWithMembership(ctx, t, instance, "IAM_ADMIN_IMPERSONATOR")
	orgUserID, orgImpersonatorPAT := createMachineUserPATWithMembership(ctx, t, instance, "ORG_ADMIN_IMPERSONATOR")
	serviceUserID, noPermPAT := createMachineUserPATWithMembership(ctx, t, instance)

	teResp, err := tokenexchange.ExchangeToken(ctx, exchanger, noPermPAT, oidc.AccessTokenType, "", "", nil, nil, nil, oidc.AccessTokenType)
	require.NoError(t, err)

	patScopes := oidc.SpaceDelimitedArray{"openid", "profile", "urn:zitadel:iam:user:metadata", "urn:zitadel:iam:user:resourceowner"}

	relyingParty, err := rp.NewRelyingPartyOIDC(ctx, instance.OIDCIssuer(), client.GetClientId(), "", "", []string{"openid"}, rp.WithJWTProfile(rp.SignerFromKeyFile(keyData)))
	require.NoError(t, err)
	resourceServer, err := instance.CreateResourceServerJWTProfile(ctx, keyData)
	require.NoError(t, err)

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
		name    string
		args    args
		want    result
		wantErr bool
	}{
		{
			name: "IMPERSONATION: subject: userID, actor: access token, membership not found error",
			args: args{
				SubjectToken:       userResp.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.AccessTokenType,
				ActorToken:         noPermPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "IMPERSONATION: subject: userID, actor: access token, requested type: JWT, membership not found error",
			args: args{
				SubjectToken:       userResp.GetUserId(),
				SubjectTokenType:   oidc_api.UserIDTokenType,
				RequestedTokenType: oidc.JWTTokenType,
				ActorToken:         noPermPAT,
				ActorTokenType:     oidc.AccessTokenType,
			},
			wantErr: true,
		},
		{
			name: "Instance IMPERSONATION: subject: userID, actor: access token, success",
			args: args{
				SubjectToken:       userResp.GetUserId(),
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, userResp.GetUserId(), iamUserID),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, userResp.GetUserId(), iamUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: userID, actor: access token, success",
			args: args{
				SubjectToken:       userResp.GetUserId(),
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, userResp.GetUserId(), orgUserID),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, userResp.GetUserId(), orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: access token, actor: access token, success",
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, serviceUserID, orgUserID),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, serviceUserID, orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: ID token, actor: access token, success",
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, serviceUserID, orgUserID),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, serviceUserID, orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: JWT, actor: access token, success",
			args: args{
				SubjectToken: func() string {
					token, err := crypto.Sign(&oidc.JWTTokenRequest{
						Issuer:    client.GetClientId(),
						Subject:   userResp.GetUserId(),
						Audience:  oidc.Audience{instance.OIDCIssuer()},
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
				verifyAccessToken: accessTokenVerifier(ctx, resourceServer, userResp.GetUserId(), orgUserID),
				verifyIDToken:     idTokenVerifier(ctx, relyingParty, userResp.GetUserId(), orgUserID),
			},
		},
		{
			name: "ORG IMPERSONATION: subject: access token, actor: access token, with refresh token, success",
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
				verifyAccessToken:  accessTokenVerifier(ctx, resourceServer, serviceUserID, orgUserID),
				verifyIDToken:      idTokenVerifier(ctx, relyingParty, serviceUserID, orgUserID),
				verifyRefreshToken: refreshTokenVerifier(ctx, relyingParty, serviceUserID, orgUserID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenexchange.ExchangeToken(ctx, exchanger, tt.args.SubjectToken, tt.args.SubjectTokenType, tt.args.ActorToken, tt.args.ActorTokenType, tt.args.Resource, tt.args.Audience, tt.args.Scopes, tt.args.RequestedTokenType)
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

// TestServer_TokenExchangeAdminImpersonation_GHSA_w4gv_rcwj_w6r5 asserts that an actor which
// only holds the `impersonation` permission cannot impersonate an administrator.
// Impersonating an administrator requires the `admin.impersonation` permission, which is
// only granted by the *_ADMIN_IMPERSONATOR roles.
// A user is considered an administrator as soon as they have any membership.
func TestServer_TokenExchangeAdminImpersonation_GHSA_w4gv_rcwj_w6r5(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ctx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	setImpersonationPolicy(t, instance, true)

	client, keyData, err := instance.CreateOIDCTokenExchangeClient(ctx, t)
	require.NoError(t, err)
	signer, err := rp.SignerFromKeyFile(keyData)()
	require.NoError(t, err)
	exchanger, err := tokenexchange.NewTokenExchangerJWTProfile(ctx, instance.OIDCIssuer(), client.GetClientId(), signer)
	require.NoError(t, err)
	relyingParty, err := rp.NewRelyingPartyOIDC(ctx, instance.OIDCIssuer(), client.GetClientId(), "", "", []string{"openid"}, rp.WithJWTProfile(rp.SignerFromKeyFile(keyData)))
	require.NoError(t, err)
	resourceServer, err := instance.CreateResourceServerJWTProfile(ctx, keyData)
	require.NoError(t, err)

	// ORG_END_USER_IMPERSONATOR is only granted the `impersonation` permission,
	// ORG_ADMIN_IMPERSONATOR is granted `admin.impersonation` on top of it.
	endUserImpersonatorID, endUserImpersonatorPAT := createMachineUserPATWithMembership(ctx, t, instance, "ORG_END_USER_IMPERSONATOR")
	adminImpersonatorID, adminImpersonatorPAT := createMachineUserPATWithMembership(ctx, t, instance, "ORG_ADMIN_IMPERSONATOR")

	// A user without any membership is a regular end user.
	regularUser := instance.CreateHumanUser(ctx)
	// A user with an org membership is an administrator.
	adminUserID, _ := createMachineUserPATWithMembership(ctx, t, instance, "ORG_OWNER")

	t.Run("end user impersonator may impersonate a regular user", func(t *testing.T) {
		resp := awaitImpersonationAllowed(ctx, t, exchanger, regularUser.GetUserId(), endUserImpersonatorPAT, oidc.AccessTokenType)
		accessTokenVerifier(ctx, resourceServer, regularUser.GetUserId(), endUserImpersonatorID)(t, resp.AccessToken)
		idTokenVerifier(ctx, relyingParty, regularUser.GetUserId(), endUserImpersonatorID)(t, resp.IDToken)
	})

	t.Run("SECURITY: end user impersonator may not impersonate an admin user", func(t *testing.T) {
		awaitImpersonationDenied(ctx, t, exchanger, adminUserID, endUserImpersonatorPAT, oidc.AccessTokenType)
	})

	t.Run("SECURITY: end user impersonator may not impersonate an admin user, requested type: JWT", func(t *testing.T) {
		awaitImpersonationDenied(ctx, t, exchanger, adminUserID, endUserImpersonatorPAT, oidc.JWTTokenType)
	})

	// The subtests above already waited for the membership of adminUserID to be projected,
	// so the following case is guaranteed to take the admin impersonation path.
	t.Run("admin impersonator may impersonate an admin user", func(t *testing.T) {
		resp := awaitImpersonationAllowed(ctx, t, exchanger, adminUserID, adminImpersonatorPAT, oidc.AccessTokenType)
		accessTokenVerifier(ctx, resourceServer, adminUserID, adminImpersonatorID)(t, resp.AccessToken)
		idTokenVerifier(ctx, relyingParty, adminUserID, adminImpersonatorID)(t, resp.IDToken)
	})

	t.Run("admin impersonator may impersonate a regular user", func(t *testing.T) {
		resp := awaitImpersonationAllowed(ctx, t, exchanger, regularUser.GetUserId(), adminImpersonatorPAT, oidc.AccessTokenType)
		accessTokenVerifier(ctx, resourceServer, regularUser.GetUserId(), adminImpersonatorID)(t, resp.AccessToken)
		idTokenVerifier(ctx, relyingParty, regularUser.GetUserId(), adminImpersonatorID)(t, resp.IDToken)
	})
}

// This test tries to call the zitadel API with an impersonated token,
// which should fail.
func TestImpersonation_API_Call(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ctx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	client, keyData, err := instance.CreateOIDCTokenExchangeClient(ctx, t)
	require.NoError(t, err)
	signer, err := rp.SignerFromKeyFile(keyData)()
	require.NoError(t, err)
	exchanger, err := tokenexchange.NewTokenExchangerJWTProfile(ctx, instance.OIDCIssuer(), client.GetClientId(), signer)
	require.NoError(t, err)
	resourceServer, err := instance.CreateResourceServerJWTProfile(ctx, keyData)
	require.NoError(t, err)

	setImpersonationPolicy(t, instance, true)

	iamUserID, iamImpersonatorPAT := createMachineUserPATWithMembership(ctx, t, instance, "IAM_ADMIN_IMPERSONATOR")
	iamOwner := instance.Users.Get(integration.UserTypeIAMOwner)

	// impersonating the instance owner!
	resp, err := tokenexchange.ExchangeToken(ctx, exchanger, iamOwner.Token, oidc.AccessTokenType, iamImpersonatorPAT, oidc.AccessTokenType, nil, nil, nil, oidc.AccessTokenType)
	require.NoError(t, err)
	accessTokenVerifier(ctx, resourceServer, iamOwner.ID, iamUserID)

	impersonatedCTX := integration.WithAuthorizationToken(ctx, resp.AccessToken)
	_, err = instance.Client.Admin.GetAllowedLanguages(impersonatedCTX, &admin.GetAllowedLanguagesRequest{})
	status := status.Convert(err)
	assert.Equal(t, codes.PermissionDenied, status.Code())
	assert.Equal(t, "Errors.TokenExchange.Token.NotForAPI (APP-Shi0J)", status.Message())
}
