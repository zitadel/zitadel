//go:build integration

package oidc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX               context.Context
	Tester            *integration.Tester
	Client            session.SessionServiceClient
	User              *user.AddHumanUserResponse
	GenericOAuthIDPID string
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
		Client = Tester.Client.SessionV2

		CTX, _ = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		return m.Run()
	}())
}

func createClient(t testing.TB) string {
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI)
	require.NoError(t, err)
	return app.GetClientId()
}

func createAuthRequest(t testing.TB, clientID string) string {
	redURL, err := Tester.CreateOIDCAuthRequest(CTX, clientID, "loginClient", redirectURI)
	require.NoError(t, err)
	return redURL
}

func TestOPStorage_CreateAuthRequest(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID)
	require.Contains(t, id, oidc_internal.IDPrefix)
}
