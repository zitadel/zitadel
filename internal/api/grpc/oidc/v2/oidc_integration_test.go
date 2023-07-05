//go:build integration

package oidc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client oidc_pb.OIDCServiceClient
	User   *user.AddHumanUserResponse
)

const (
	redirectURI = "oidc_integration_test://callback"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.OIDCv2

		CTX = Tester.WithSystemAuthorization(ctx, integration.OrgOwner)
		User = Tester.CreateHumanUser(CTX)
		return m.Run()
	}())
}

func TestServer_GetAuthRequest(t *testing.T) {
	client, err := Tester.CreateOIDCNativeClient(CTX, redirectURI)
	require.NoError(t, err)
	authRequestID, err := Tester.CreateOIDCAuthRequest(client.GetClientId(), Tester.Users[integration.OrgOwner].ID, redirectURI)
	require.NoError(t, err)
	now := time.Now()

	tests := []struct {
		name          string
		AuthRequestID string
		want          *oidc_pb.GetAuthRequestResponse
		wantErr       bool
	}{
		{
			name:          "Not found",
			AuthRequestID: "123",
			wantErr:       true,
		},
		{
			name:          "success",
			AuthRequestID: authRequestID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.GetAuthRequest(CTX, &oidc_pb.GetAuthRequestRequest{
				AuthRequestId: tt.AuthRequestID,
			})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			authRequest := got.GetAuthRequest()
			assert.NotNil(t, authRequest)
			assert.Equal(t, authRequestID, authRequest.GetId())
			assert.WithinRange(t, authRequest.GetCreationDate().AsTime(), now.Add(-time.Second), now.Add(time.Second))
			assert.Contains(t, authRequest.GetScope(), "openid")
		})
	}
}

func TestServer_LinkSessionToAuthRequest(t *testing.T) {
	client, err := Tester.CreateOIDCNativeClient(CTX, redirectURI)
	require.NoError(t, err)
	authRequestID, err := Tester.CreateOIDCAuthRequest(client.GetClientId(), integration.LoginUser, redirectURI)
	require.NoError(t, err)
	sessionResp, err := Tester.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: Tester.Users[integration.OrgOwner].ID,
				},
			},
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name          string
		AuthRequestID string
		SessionID     string
		SessionToken  string
		AuthError     string
		want          *oidc_pb.LinkSessionToAuthRequestResponse
		wantErr       bool
	}{
		{
			name:          "Not found",
			AuthRequestID: "123",
			SessionID:     "",
			SessionToken:  "",
			wantErr:       true,
		},
		{
			name:          "auth error",
			AuthRequestID: "",
			SessionID:     "",
			SessionToken:  "",
			wantErr:       true,
		},
		{
			name:          "session not found",
			AuthRequestID: authRequestID,
			SessionID:     "",
			SessionToken:  "",
			wantErr:       true,
		},
		{
			name:          "session token invalid",
			AuthRequestID: authRequestID,
			SessionID:     sessionResp.SessionId,
			SessionToken:  "invalid",
			wantErr:       true,
		},
		{
			name:          "success",
			AuthRequestID: authRequestID,
			SessionID:     sessionResp.SessionId,
			SessionToken:  sessionResp.SessionToken,
			want: &oidc_pb.LinkSessionToAuthRequestResponse{
				CallbackUrl: "/oauth/v2/authorize/callback?id=" + authRequestID,
				Details: &object.Details{
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.LinkSessionToAuthRequest(CTX, &oidc_pb.LinkSessionToAuthRequestRequest{
				AuthRequestId: tt.AuthRequestID,
				SessionId:     tt.SessionID,
				SessionToken:  tt.SessionToken,
			})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
