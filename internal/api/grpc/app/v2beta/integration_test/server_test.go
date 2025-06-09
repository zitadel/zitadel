//go:build integration

package instance_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
)

var (
	CTX                  context.Context
	AdminCTX             context.Context
	instance             *integration.Instance
	instancePermissionV2 *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		CTX = ctx
		instance = integration.NewInstance(ctx)
		AdminCTX = instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)

		return m.Run()
	}())
}
