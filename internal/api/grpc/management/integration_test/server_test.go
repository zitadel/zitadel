//go:build integration

package management_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/v2/internal/integration"
	mgmt_pb "github.com/zitadel/zitadel/v2/pkg/grpc/management"
)

var (
	CTX, OrgCTX context.Context
	Instance    *integration.Instance
	Client      mgmt_pb.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		CTX = ctx
		OrgCTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		Client = Instance.Client.Mgmt
		return m.Run()
	}())
}
