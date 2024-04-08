//go:build integration

package management_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

var (
	CTX, OrgCTX context.Context
	Tester      *integration.Tester
	Client      mgmt_pb.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(3 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX = ctx
		OrgCTX = Tester.WithAuthorization(ctx, integration.OrgOwner)
		Client = Tester.Client.Mgmt
		return m.Run()
	}())
}
