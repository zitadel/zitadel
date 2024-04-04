//go:build integration

package oidc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/integration"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
)

func TestServer_UserInfo(t *testing.T) {
	iamOwnerCTX := Tester.WithAuthorization(CTX, integration.IAMOwner)
	t.Cleanup(func() {
		_, err := Tester.Client.FeatureV2.ResetInstanceFeatures(iamOwnerCTX, &feature.ResetInstanceFeaturesRequest{})
		require.NoError(t, err)
	})
	tests := []struct {
		name    string
		legacy  bool
		trigger bool
	}{
		{
			name:   "legacy enabled",
			legacy: true,
		},
		{
			name:    "legacy and trigger disabled",
			legacy:  false,
			trigger: false,
		},
		{
			name:    "legacy disabled, trigger enabled",
			legacy:  false,
			trigger: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Tester.Client.FeatureV2.SetInstanceFeatures(iamOwnerCTX, &feature.SetInstanceFeaturesRequest{
				OidcLegacyIntrospection:             &tt.legacy,
				OidcTriggerIntrospectionProjections: &tt.trigger,
			})
			require.NoError(t, err)
			testServer_UserInfo(t)
		})
	}
}

func testServer_UserInfo(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime)

	// test actual userinfo
	provider, err := Tester.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)
	userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
	require.NoError(t, err)
	assertUserinfo(t, userinfo)
}

func assertUserinfo(t *testing.T, userinfo *oidc.UserInfo) {
	assert.Equal(t, User.GetUserId(), userinfo.Subject)
	assert.Equal(t, "Mickey", userinfo.GivenName)
	assert.Equal(t, "Mouse", userinfo.FamilyName)
	assert.Equal(t, "Mickey Mouse", userinfo.Name)
	assert.NotEmpty(t, userinfo.PreferredUsername)
	assert.Equal(t, userinfo.PreferredUsername, userinfo.Email)
	assert.False(t, bool(userinfo.EmailVerified))
	assertOIDCTime(t, userinfo.UpdatedAt, User.GetDetails().GetChangeDate().AsTime())
}
