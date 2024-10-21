//go:build integration

package action_test

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
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
)

var (
	CTX context.Context
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		CTX = ctx
		return m.Run()
	}())
}

func ensureFeatureEnabled(t *testing.T, instance *integration.Instance) {
	ctx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.Actions.GetEnabled() {
		return
	}
	_, err = instance.Client.FeatureV2.SetInstanceFeatures(ctx, &feature.SetInstanceFeaturesRequest{
		Actions: gu.Ptr(true),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			assert.NoError(ttt, err)
			assert.True(ttt, f.Actions.GetEnabled())
		},
		retryDuration,
		tick,
		"timed out waiting for ensuring instance feature")

	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			_, err := instance.Client.ActionV3Alpha.ListExecutionMethods(ctx, &action.ListExecutionMethodsRequest{})
			assert.NoError(ttt, err)
		},
		retryDuration,
		tick,
		"timed out waiting for ensuring instance feature call")
}
