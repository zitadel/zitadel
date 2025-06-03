//go:build integration

package authorization_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"os"
	"testing"
	"time"
)

var (
	EmptyCTX         context.Context
	IAMCTX           context.Context
	NoPermissionsCTX context.Context
	Instance         *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		EmptyCTX = ctx
		Instance = integration.NewInstance(ctx)
		IAMCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		NoPermissionsCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeNoPermission)
		return m.Run()
	}())
}

func setPermissionCheckV2Flag(t *testing.T, enabled bool) {
	integration.EnsureInstanceFeature(t, IAMCTX, Instance, &feature.SetInstanceFeaturesRequest{
		PermissionCheckV2: &enabled,
	}, func(tt *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.Equal(tt, enabled, got.GetPermissionCheckV2().GetEnabled(), "permission check v2 feature should be set to %v", enabled)
	})
}
