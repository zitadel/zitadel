//go:build integration

package oidc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/api/oidc/amr"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/integration"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	User   *user.AddHumanUserResponse
)

const (
	redirectURI = "oidc_integration_test://callback"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, _ = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		User = Tester.CreateHumanUser(CTX)
		Tester.RegisterUserPasskey(CTX, User.GetUserId())
		return m.Run()
	}())
}

func createClient(t testing.TB) string {
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI)
	require.NoError(t, err)
	return app.GetClientId()
}

func createAuthRequest(t testing.TB, clientID string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequest(clientID, "loginClient", redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func TestOPStorage_CreateAuthRequest(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID)
	require.Contains(t, id, command.IDPrefixV2)
}

func TestOPStorage_CreateAccessToken(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID)
	_ = id
	createResp, err := Tester.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTX, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	// link session to auth request and get code
	var code string

	tokens := token(t, clientID, code)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.IDToken)
	assert.Empty(t, tokens.RefreshToken)
	assert.Equal(t, []string{amr.UserPresence, amr.MFA}, tokens.IDTokenClaims.AuthenticationMethodsReferences)
	assert.Equal(t, updateResp.Details.ChangeDate, tokens.IDTokenClaims.AuthTime)
}

func TestOPStorage_CreateAccessAndRefreshTokens(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	_ = id
	createResp, err := Tester.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTX, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	// link session to auth request and get code
	var code string

	tokens := token(t, clientID, code)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.IDToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, []string{amr.UserPresence, amr.MFA}, tokens.IDTokenClaims.AuthenticationMethodsReferences)
	assert.Equal(t, updateResp.Details.ChangeDate, tokens.IDTokenClaims.AuthTime)
}

func TestOPStorage_CreateAccessAndRefreshTokens_refresh(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	_ = id
	createResp, err := Tester.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTX, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	// link session to auth request and get code
	_ = updateResp
	var code string

	tokens := token(t, clientID, code)
	assert.NotEmpty(t, tokens.RefreshToken)

	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	newTokens, err := refreshToken(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.Extra("id_token"))
	assert.NotEmpty(t, newTokens.RefreshToken)

	// refresh with an old refresh_token must fail
	_, err = rp.RefreshAccessToken(provider, tokens.RefreshToken, "", "")
	require.Error(t, err)
}

func token(t testing.TB, clientID, code string) *oidc.Tokens[*oidc.IDTokenClaims] {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	codeVerifier := "codeVerifier"
	tokens, err := rp.CodeExchange[*oidc.IDTokenClaims](context.Background(), code, provider, rp.WithCodeVerifier(codeVerifier))
	require.NoError(t, err)

	return tokens
}

func refreshToken(t testing.TB, clientID, refreshToken string) (*oauth2.Token, error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	return rp.RefreshAccessToken(provider, refreshToken, "", "")
}
