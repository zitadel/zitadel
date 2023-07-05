//go:build integration

package handlers_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

var (
	SystemCTX    context.Context
	Tester       *integration.Tester
	SystemClient system.SystemServiceClient
	MgmtClient   management.ManagementServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()
		Tester = integration.NewTester(ctx)
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		MgmtClient = Tester.Client.Mgmt
		SystemClient = Tester.Client.System
		defer Tester.Done()
		return m.Run()
	}())
}
