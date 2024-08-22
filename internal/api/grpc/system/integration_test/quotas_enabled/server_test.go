//go:build integration

package quotas_enabled_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	CTX      context.Context
	Instance *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		CTX = ctx

		Instance = integration.GetInstance(ctx)
		return m.Run()
	}())
}
