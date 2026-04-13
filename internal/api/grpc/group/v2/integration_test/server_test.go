//go:build integration

package group_test

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
	CTX                  context.Context
	instance             *integration.Instance
	instancePermissionV2 *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		CTX = ctx
		instance = integration.NewInstance(ctx)
		instancePermissionV2 = integration.NewInstance(ctx)
		return m.Run()
	}())
}

func ensureFeaturePermissionV2Enabled(t *testing.T, testInstance *integration.Instance) {
	ctx := testInstance.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	f, err := testInstance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)

	if f.PermissionCheckV2.GetEnabled() {
		return
	}

	_, err = testInstance.Client.FeatureV2.SetInstanceFeatures(ctx, &feature.SetInstanceFeaturesRequest{
		PermissionCheckV2: gu.Ptr(true),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		f, err := testInstance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{Inheritance: true})
		require.NoError(tt, err)
		assert.True(tt, f.PermissionCheckV2.GetEnabled())
	}, retryDuration, tick, "timed out waiting for ensuring testInstance feature")
}
