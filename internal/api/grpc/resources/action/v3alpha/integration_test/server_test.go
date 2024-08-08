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
	feature "github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client action.ZITADELActionsClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		var err error
		Tester, err = integration.NewTester(ctx)
		if err != nil {
			panic(err)
		}
		Client = Tester.Client.ActionV3

		CTX = Tester.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		return m.Run()
	}())
}

func ensureFeatureEnabled(t *testing.T) {
	f, err := Tester.Client.FeatureV2.GetInstanceFeatures(CTX, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.Actions.GetEnabled() {
		return
	}
	_, err = Tester.Client.FeatureV2.SetInstanceFeatures(CTX, &feature.SetInstanceFeaturesRequest{
		Actions: gu.Ptr(true),
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
			if f.Actions.GetEnabled() {
				return
			}
		},
		retryDuration,
		100*time.Millisecond,
		"timed out waiting for ensuring instance feature")
}
