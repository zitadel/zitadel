//go:build integration

package oidc_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

var (
	CTX            context.Context
	CTXLoginClient context.Context
	IAMCTX         context.Context
	Instance       *integration.Instance
	Client         oidc_pb.OIDCServiceClient
	loginV2        = &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: nil}}}
)

const (
	redirectURI         = "oidcintegrationtest://callback"
	redirectURIImplicit = "http://localhost:9999/callback"
	logoutRedirectURI   = "oidcintegrationtest://logged-out"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OIDCv2

		IAMCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		CTXLoginClient = Instance.WithAuthorization(ctx, integration.UserTypeLogin)
		return m.Run()
	}())
}
