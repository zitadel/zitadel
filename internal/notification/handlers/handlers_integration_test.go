//go:build integration

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
	Tester    *integration.Tester
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()
		CTX = ctx

		Tester = integration.NewTester(ctx, `
Quotas:
  Access:
    Enabled: true
`)
		defer Tester.Done()

		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		return m.Run()
	}())
}
