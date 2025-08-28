//go:build integration

package settings_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

var (
	CTX, AdminCTX, UserTypeLoginCtx, OrgOwnerCtx context.Context
	Instance                                     *integration.Instance
	Client                                       settings.SettingsServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = ctx
		AdminCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		UserTypeLoginCtx = Instance.WithAuthorizationToken(ctx, integration.UserTypeLogin)
		OrgOwnerCtx = Instance.WithAuthorizationToken(ctx, integration.UserTypeOrgOwner)

		Client = Instance.Client.SettingsV2
		return m.Run()
	}())
}
