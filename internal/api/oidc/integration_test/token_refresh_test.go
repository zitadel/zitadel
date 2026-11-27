//go:build integration

package oidc_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc/metadata"

	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// Test_RefreshToken_failedAction checks that a refresh token request which fails in an action
// with interrupt on error does not consume the refresh token: the previous access token must stay valid,
// and the client must be able to retry with the same refresh token once the target is available again.
func Test_RefreshToken_failedAction(t *testing.T) {
	tests := []struct {
		name     string
		function string
	}{
		{
			name:     "preaccesstoken (JWT access token)",
			function: "preaccesstoken",
		},
		{
			name:     "preuserinfo (ID token)",
			function: "preuserinfo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Executions apply to the whole instance, so use a separate one.
			instance := integration.NewInstance(CTX)
			ctx := instance.WithAuthorization(CTX, integration.UserTypeOrgOwner)
			iamCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
			loginCtx := instance.WithAuthorization(CTX, integration.UserTypeLogin)

			user := instance.CreateHumanUser(ctx)
			instance.SetUserPassword(ctx, user.GetUserId(), integration.UserPassword, false)

			clientID, _ := createClient(t, instance)
			authRequestID := createAuthRequest(t, instance, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
			sessionID, sessionToken, _, _ := instance.CreatePasswordSession(t, loginCtx, user.GetUserId(), integration.UserPassword)
			linkResp, err := instance.Client.OIDCv2.CreateCallback(loginCtx, &oidc_pb.CreateCallbackRequest{
				AuthRequestId: authRequestID,
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    sessionID,
						SessionToken: sessionToken,
					},
				},
			})
			require.NoError(t, err)
			code := assertCodeResponse(t, linkResp.GetCallbackUrl())
			tokens, err := exchangeTokens(t, instance, clientID, code, redirectURI)
			require.NoError(t, err)
			assertTokens(t, tokens, true)
			require.True(t, assertAccessTokenValid(t, instance, tokens.TokenType, tokens.AccessToken, user.GetUserId()))

			// The target answers 503 until targetAvailable is set, and counts the successful calls.
			var (
				targetAvailable atomic.Bool
				targetSucceeded atomic.Int32
			)
			target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if !targetAvailable.Load() {
					http.Error(w, "unavailable", http.StatusServiceUnavailable)
					return
				}
				targetSucceeded.Add(1)
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, "{}")
			}))
			defer target.Close()
			targetResp := instance.CreateTarget(iamCtx, t, "", target.URL, target_domain.TargetTypeCall, true, action.PayloadType_PAYLOAD_TYPE_JSON)
			setFunctionExecution(iamCtx, t, instance, tt.function, targetResp.GetId())

			// The refresh fails because the target fails and interrupts the request.
			_, err = refreshTokensForInstance(t, instance, clientID, tokens.RefreshToken)
			require.ErrorContains(t, err, "Errors.Execution.Failed")

			// Nothing was issued, so the previous access token must still be valid.
			assertAccessTokenValid(t, instance, tokens.TokenType, tokens.AccessToken, user.GetUserId())

			// Once the target answers again, the same refresh token can be used.
			targetAvailable.Store(true)
			newTokens, err := refreshTokensForInstance(t, instance, clientID, tokens.RefreshToken)
			require.NoError(t, err)
			assert.NotEmpty(t, newTokens.AccessToken)
			assert.NotEmpty(t, newTokens.RefreshToken)
			assert.NotEqual(t, tokens.RefreshToken, newTokens.RefreshToken)
			assertAccessTokenValid(t, instance, newTokens.TokenType, newTokens.AccessToken, user.GetUserId())
			// The action runs once per refresh, even though it now runs before the session is updated.
			assert.EqualValues(t, 1, targetSucceeded.Load(), "target calls for the successful refresh")

			// The successful refresh renewed the refresh token, so the old one can't be used again.
			_, err = refreshTokensForInstance(t, instance, clientID, tokens.RefreshToken)
			require.ErrorContains(t, err, "Errors.OIDCSession.RefreshTokenInvalid")
		})
	}
}

func refreshTokensForInstance(t testing.TB, instance *integration.Instance, clientID, refreshToken string) (*oidc.Tokens[*oidc.IDTokenClaims], error) {
	provider, err := instance.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)
	return rp.RefreshTokens[*oidc.IDTokenClaims](CTX, provider, refreshToken, "", "")
}

func assertAccessTokenValid(t *testing.T, instance *integration.Instance, tokenType, accessToken, userID string) bool {
	t.Helper()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokenType, accessToken))
	myUserResp, err := instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	return assert.NoError(t, err) && assert.Equal(t, userID, myUserResp.GetUser().GetId())
}

// setFunctionExecution sets the targets of a function execution and waits until the execution can be listed.
func setFunctionExecution(ctx context.Context, t *testing.T, instance *integration.Instance, function string, targetIDs ...string) {
	condition := &action.Condition{
		ConditionType: &action.Condition_Function{
			Function: &action.FunctionExecution{Name: function},
		},
	}
	instance.SetExecution(ctx, t, condition, targetIDs)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, err := instance.Client.ActionV2.ListExecutions(ctx, &action.ListExecutionsRequest{
			Filters: []*action.ExecutionSearchFilter{
				{Filter: &action.ExecutionSearchFilter_InConditionsFilter{
					InConditionsFilter: &action.InConditionsFilter{Conditions: []*action.Condition{condition}},
				}},
			},
		})
		if !assert.NoError(ttt, err) {
			return
		}
		if assert.Len(ttt, got.GetExecutions(), 1) {
			assert.Len(ttt, got.GetExecutions()[0].GetTargets(), len(targetIDs))
		}
	}, retryDuration, tick, "timeout waiting for the execution")
}
