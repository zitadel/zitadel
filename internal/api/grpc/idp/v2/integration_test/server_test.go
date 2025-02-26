//go:build integration

package idp_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
)

var (
	CTX      context.Context
	IamCTX   context.Context
	UserCTX  context.Context
	Instance *integration.Instance
	Client   idp_pb.IdentityProviderServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		UserCTX = Instance.WithAuthorization(ctx, integration.UserTypeNoPermission)
		IamCTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		Client = Instance.Client.IDPv2
		return m.Run()
	}())
}
