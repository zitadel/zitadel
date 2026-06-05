//go:build integration

package authorization_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	EmptyCTX             context.Context
	IAMCTX               context.Context
	Instance             *integration.Instance
	InstancePermissionV2 *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		EmptyCTX = ctx
		Instance = integration.NewInstance(ctx)
		IAMCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		InstancePermissionV2 = integration.NewInstance(ctx)
		return m.Run()
	}())
}
