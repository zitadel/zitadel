//go:build integration

package oidc_test

import (
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

func TestServer_DeviceAuth(t *testing.T) {
	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	client, err := Instance.CreateOIDCClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, app.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE)
	require.NoError(t, err)

	tests := []struct {
		name     string
		scope    []string
		decision func(t *testing.T, id string)
		wantErr  error
	}{
		{
			name:  "authorized",
			scope: []string{},
			decision: func(t *testing.T, id string) {
				sessionID, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
				_, err = Instance.Client.OIDCv2.AuthorizeOrDenyDeviceAuthorization(CTXLOGIN, &oidc_pb.AuthorizeOrDenyDeviceAuthorizationRequest{
					DeviceAuthorizationId: id,
					Decision: &oidc_pb.AuthorizeOrDenyDeviceAuthorizationRequest_AuthorizeWithSession{
						AuthorizeWithSession: &oidc_pb.Session{
							SessionId:    sessionID,
							SessionToken: sessionToken,
						},
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name:  "authorized, with ZITADEL",
			scope: []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, domain.ProjectScopeZITADEL},
			decision: func(t *testing.T, id string) {
				sessionID, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
				_, err = Instance.Client.OIDCv2.AuthorizeOrDenyDeviceAuthorization(CTXLOGIN, &oidc_pb.AuthorizeOrDenyDeviceAuthorizationRequest{
					DeviceAuthorizationId: id,
					Decision: &oidc_pb.AuthorizeOrDenyDeviceAuthorizationRequest_AuthorizeWithSession{
						AuthorizeWithSession: &oidc_pb.Session{
							SessionId:    sessionID,
							SessionToken: sessionToken,
						},
					},
				})
				require.NoError(t, err)
			},
		},
		{
			name:  "denied",
			scope: []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, domain.ProjectScopeZITADEL},
			decision: func(t *testing.T, id string) {
				_, err = Instance.Client.OIDCv2.AuthorizeOrDenyDeviceAuthorization(CTXLOGIN, &oidc_pb.AuthorizeOrDenyDeviceAuthorizationRequest{
					DeviceAuthorizationId: id,
					Decision: &oidc_pb.AuthorizeOrDenyDeviceAuthorizationRequest_Deny{
						Deny: &oidc_pb.Deny{},
					},
				})
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			provider, err := rp.NewRelyingPartyOIDC(CTX, Instance.OIDCIssuer(), client.GetClientId(), "", "", tt.scope)
			require.NoError(t, err)
			deviceAuthorization, err := rp.DeviceAuthorization(CTX, tt.scope, provider, nil)
			require.NoError(t, err)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			var tokens *oidc.AccessTokenResponse
			go func() {
				tokens, err = rp.DeviceAccessToken(CTX, deviceAuthorization.DeviceCode, time.Duration(deviceAuthorization.Interval)*time.Second, provider)
				require.ErrorIs(t, err, tt.wantErr)
				wg.Done()
			}()

			req, err := Instance.Client.OIDCv2.GetDeviceAuthorizationRequest(CTX, &oidc_pb.GetDeviceAuthorizationRequestRequest{
				UserCode: deviceAuthorization.UserCode,
			})
			require.NoError(t, err)

			tt.decision(t, req.GetDeviceAuthorizationRequest().GetId())

			// wait on the token response
			wg.Wait()

			_, err = Instance.Client.Auth.GetMyUser(integration.WithAuthorizationToken(CTX, tokens.AccessToken), &auth.GetMyUserRequest{})
			if slices.Contains(tt.scope, domain.ProjectScopeZITADEL) {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
