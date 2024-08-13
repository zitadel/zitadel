//go:build integration_old

package handlers_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	CTX       context.Context
	SystemCTX context.Context
	Instance  *integration.Instance
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
		defer Instance.Done()

		SystemCTX = Instance.WithAuthorization(ctx, integration.SystemUser)
		return m.Run()
	}())
}
