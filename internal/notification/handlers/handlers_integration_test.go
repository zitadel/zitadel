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
		CTX = ctx
		defer cancel()
		Tester = integration.NewTester(ctx)
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		defer Tester.Done()
		return m.Run()
	}())
}
