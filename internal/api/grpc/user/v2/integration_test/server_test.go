//go:build integration

package user_test

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
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX                            context.Context
	OrgCTX                         context.Context
	IamCTX                         context.Context
	LoginCTX                       context.Context
	UserCTX                        context.Context
	SystemCTX                      context.Context
	SystemUserWithNoPermissionsCTX context.Context
	Instance                       *integration.Instance
	InstancePermissionV2           *integration.Instance
	Client                         user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		CTX = ctx

		Instance = integration.NewInstance(ctx)
		InstancePermissionV2 = integration.NewInstance(ctx)

		SystemUserWithNoPermissionsCTX = integration.WithSystemUserWithNoPermissionsAuthorization(ctx)
		UserCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeNoPermission)
		IamCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		LoginCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeLogin)
		SystemCTX = integration.WithSystemAuthorization(ctx)
		OrgCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeOrgOwner)
		Client = Instance.Client.UserV2
		return m.Run()
	}())
}

func ensureFeaturePermissionV2Enabled(t *testing.T, instance *integration.Instance) {
	ctx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
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
