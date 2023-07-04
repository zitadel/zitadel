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
	CTX          context.Context
	Tester       *integration.Tester
	SystemClient system.SystemServiceClient
	MgmtClient   management.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		CTX = ctx
		defer cancel()
		os.Setenv("INTEGRATION_DB_FLAVOR", "postgres")
		os.Setenv("ZITADEL_MASTERKEY", "MasterkeyNeedsToHave32Characters")
		Tester = integration.NewTester(ctx)
		MgmtClient = Tester.Client.Mgmt
		SystemClient = Tester.Client.System
		defer Tester.Done()
		return m.Run()
	}())
}
