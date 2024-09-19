//go:build integration

package userschema_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/v2/internal/integration"
	"github.com/zitadel/zitadel/v2/pkg/grpc/feature/v2"
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
	if f.UserSchema.GetEnabled() {
		return
	}
	_, err = instance.Client.FeatureV2.SetInstanceFeatures(ctx, &feature.SetInstanceFeaturesRequest{
		UserSchema: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration := time.Minute
	if ctxDeadline, ok := ctx.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			require.NoError(ttt, err)
			if f.UserSchema.GetEnabled() {
				return
			}
		},
		retryDuration,
		time.Second,
		"timed out waiting for ensuring instance feature")
}
