//go:build integration_old

package quotas_enabled_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	CTX         context.Context
	SystemCTX   context.Context
	IAMOwnerCTX context.Context
	Instance    *integration.Tester
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()
		CTX = ctx

		Instance = integration.NewTester(ctx, `
Quotas:
  Access:
    Enabled: true
`)
		defer Tester.Done()
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		IAMOwnerCTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		return m.Run()
	}())
}
