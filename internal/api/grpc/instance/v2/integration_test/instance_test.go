//go:build integration

package instance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func TestDeleteInstance(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	// Event store instances
	instES := integration.NewInstance(ctxWithSysAuthZ)
	instanceESOwnerCtx := instES.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)
	instES2 := integration.NewInstance(ctxWithSysAuthZ)

	instRelational := integration.NewInstance(ctxWithSysAuthZ)
	integration.EnsureInstanceFeature(t, ctxWithSysAuthZ, instRelational, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})
	instanceRelationalOwnerCtx := instRelational.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)
	instRelational2 := integration.NewInstance(ctxWithSysAuthZ)
	integration.EnsureInstanceFeature(t, ctxWithSysAuthZ, instRelational2, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})

	type instAndCtx struct {
		testType string
		i1       *integration.Instance
		i2       *integration.Instance
		ctx      context.Context
	}

	testedInstancesAndCtxs := []instAndCtx{
		{testType: "eventstore", i1: instES, i2: instES2, ctx: instanceESOwnerCtx},
		{testType: "relational", i1: instRelational, i2: instRelational2, ctx: instanceRelationalOwnerCtx},
	}

	t.Cleanup(func() {
		instES.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instES.ID()})
		instES2.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instES2.ID()})
		instRelational.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instRelational.ID()})
		instRelational2.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instRelational2.ID()})
	})

	for _, instanceWithCtx := range testedInstancesAndCtxs {
		tt := []struct {
			testName             string
			inputRequest         *instance.DeleteInstanceRequest
			inputContext         context.Context
			expectedErrorMsg     string
			expectedErrorCode    codes.Code
			expectedInstanceID   string
			expectsNullTimestamp bool
		}{
			{
				testName: "when invalid context should return unauthN error",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i1.ID(),
				},
				inputContext:      t.Context(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				testName: "when invalid context and invalid id should return unauthN error",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i1.ID() + "invalid",
				},
				inputContext:      t.Context(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				testName: "when instance token for invalid instance should return unauthN error",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i1.ID() + "invalid",
				},
				inputContext:      instanceWithCtx.ctx,
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
			},
			{
				testName: "when instance token for other instance should return unauthN error",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i2.ID(),
				},
				inputContext:      instanceWithCtx.ctx,
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
			},
			{
				testName: "when instance token for own instance should return unauthZ error",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i1.ID(),
				},
				inputContext:      instanceWithCtx.ctx,
				expectedErrorCode: codes.PermissionDenied,
				expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
			},
			{
				testName: "when invalid id should return no error",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i1.ID() + "invalid",
				},
				inputContext:         ctxWithSysAuthZ,
				expectsNullTimestamp: true,
			},
			{
				testName: "when delete succeeds should return deletion date",
				inputRequest: &instance.DeleteInstanceRequest{
					InstanceId: instanceWithCtx.i1.ID(),
				},
				inputContext:       ctxWithSysAuthZ,
				expectedInstanceID: instanceWithCtx.i1.ID(),
			},
		}

		for _, tc := range tt {
			tName := fmt.Sprintf("%s - %s", instanceWithCtx.testType, tc.testName)
			t.Run(tName, func(t *testing.T) {
				// Test
				res, err := instanceWithCtx.i1.Client.InstanceV2.DeleteInstance(tc.inputContext, tc.inputRequest)

				// Verify
				assert.Equal(t, tc.expectedErrorCode, status.Code(err))
				assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())
				if tc.expectedErrorMsg == "" {
					require.NotNil(t, res)
					require.Equal(t, tc.expectsNullTimestamp, res.GetDeletionDate() == nil)
				}
			})
		}

	}

}

func TestUpdateInstance(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	instanceOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	inst2 := integration.NewInstance(ctxWithSysAuthZ)

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName          string
		inputRequest      *instance.UpdateInstanceRequest
		inputContext      context.Context
		expectedErrorMsg  string
		expectedErrorCode codes.Code
		expectedNewName   string
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID(),
				InstanceName: " ",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when invalid context and invalid id should return unauthN error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID() + "invalid",
				InstanceName: " ",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when instance token for invalid instance should return unauthN error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID() + "invalid",
				InstanceName: " ",
			},
			inputContext:      instanceOwnerCtx,
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
		},
		{
			testName: "when instance token for other instance should return unauthN error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst2.ID(),
				InstanceName: " ",
			},
			inputContext:      instanceOwnerCtx,
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID(),
				InstanceName: " ",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			testName: "when invalid id should return not found error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID() + "invalid",
				InstanceName: "an-updated-name",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "Errors.Instance.NotFound (INST-nuso2m)",
		},
		{
			testName: "when update succeeds (instance owner) should change instance name",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID(),
				InstanceName: "an-updated-name",
			},
			inputContext:    instanceOwnerCtx,
			expectedNewName: "an-updated-name",
		},
		{
			testName: "when update succeeds should change instance name",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst2.ID(),
				InstanceName: "an-updated-name2",
			},
			inputContext:    ctxWithSysAuthZ,
			expectedNewName: "an-updated-name2",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2.UpdateInstance(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())
			if tc.expectedErrorMsg == "" {

				require.NotNil(t, res)
				assert.NotEmpty(t, res.GetChangeDate())

				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputContext, time.Minute)
				require.EventuallyWithT(t, func(tt *assert.CollectT) {
					retrievedInstance, err := inst.Client.InstanceV2.GetInstance(tc.inputContext, &instance.GetInstanceRequest{InstanceId: tc.inputRequest.GetInstanceId()})
					require.Nil(tt, err)
					assert.Equal(tt, tc.expectedNewName, retrievedInstance.GetInstance().GetName())
				}, retryDuration, tick, "timeout waiting for expected execution result")
			}
		})
	}
}
