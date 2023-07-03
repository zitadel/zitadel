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
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
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
	authRequestID, err := Tester.CreateOIDCAuthRequest(CTX, client.GetClientId(), Tester.Users[integration.OrgOwner].ID, redirectURI)
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
