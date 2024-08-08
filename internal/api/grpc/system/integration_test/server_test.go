//go:build integration

package system_test

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
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Tester, err = integration.NewTester(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		SystemCTX = Tester.WithAuthorization(ctx, integration.UserTypeSystem)
		return m.Run()
	}())
}
