//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	Instance *integration.Instance
	CTX      context.Context
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		return m.Run()
	}())
}
