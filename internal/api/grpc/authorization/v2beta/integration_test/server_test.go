//go:build integration

package authorization_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
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

func ensureFeaturePermissionV2Enabled(t *testing.T, instance *integration.Instance) {
	ctx := instance.WithAuthorizationToken(EmptyCTX, integration.UserTypeIAMOwner)
	f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.PermissionCheckV2.GetEnabled() {
		return
	}
	_, err = instance.Client.FeatureV2.SetInstanceFeatures(ctx, &feature.SetInstanceFeaturesRequest{
		PermissionCheckV2: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			assert.NoError(ttt, err)
			if f.PermissionCheckV2.GetEnabled() {
				return
			}
		},
		retryDuration,
		tick,
		"timed out waiting for ensuring instance feature")
}
