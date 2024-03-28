//go:build integration

package execution_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2beta"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client execution.ExecutionServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.ExecutionV3

		CTX, _ = Tester.WithAuthorization(ctx, integration.IAMOwner), errCtx
		return m.Run()
	}())
}

func ensureFeatureEnabled(t *testing.T) {
	f, err := Tester.Client.FeatureV2.GetInstanceFeatures(CTX, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.Execution.GetEnabled() {
		return
	}
	_, err = Tester.Client.FeatureV2.SetInstanceFeatures(CTX, &feature.SetInstanceFeaturesRequest{
		Execution: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration := time.Minute
	if ctxDeadline, ok := CTX.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := Tester.Client.FeatureV2.GetInstanceFeatures(CTX, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			require.NoError(ttt, err)
			if f.Execution.GetEnabled() {
				return
			}
		},
		retryDuration,
		100*time.Millisecond,
		"timed out waiting for ensuring instance feature")
}
