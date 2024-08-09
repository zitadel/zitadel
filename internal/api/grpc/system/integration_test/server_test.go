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
	Instance  *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Instance, err = integration.FirstInstance(ctx)
		if err != nil {
			panic(err)
		}

		CTX = ctx
		SystemCTX = Instance.WithAuthorization(ctx, integration.UserTypeSystem)
		return m.Run()
	}())
}
