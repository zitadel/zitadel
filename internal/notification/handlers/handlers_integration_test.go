//go:build integration

package handlers_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/pkg/grpc/management"

	"github.com/zitadel/zitadel/pkg/grpc/system"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	Tester                    *integration.Tester
	SystemUserCTX             context.Context
	SystemClient              system.SystemServiceClient
	PrimaryDomain, InstanceID string
	IAMOwnerCtx               context.Context
	MgmtClient                management.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()
		os.Setenv("INTEGRATION_DB_FLAVOR", "postgres")
		os.Setenv("ZITADEL_MASTERKEY", "MasterkeyNeedsToHave32Characters")
		Tester = integration.NewTester(ctx)
		PrimaryDomain, InstanceID, SystemUserCTX, IAMOwnerCtx = Tester.UseIsolatedInstance(ctx)
		MgmtClient = Tester.Client.Mgmt
		SystemClient = Tester.Client.System
		defer Tester.Done()
		return m.Run()
	}())
}
