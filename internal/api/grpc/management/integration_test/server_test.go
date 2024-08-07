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
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Tester, err = integration.NewTester(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		OrgCTX = Tester.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		Client = Tester.Client.Mgmt
		return m.Run()
	}())
}
