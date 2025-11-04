//go:build integration

package instance_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func TestDeleteInstance(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	instanceOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	inst2 := integration.NewInstance(ctxWithSysAuthZ)

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName           string
		inputRequest       *instance.DeleteInstanceRequest
		inputContext       context.Context
		expectedErrorMsg   string
		expectedErrorCode  codes.Code
		expectedInstanceID string
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID(),
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when invalid context and invalid id should return unauthN error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID() + "invalid",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when instance token for invalid instance should return unauthN error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID() + "invalid",
			},
			inputContext:      instanceOwnerCtx,
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
		},
		{
			testName: "when instance token for other instance should return unauthN error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst2.ID(),
			},
			inputContext:      instanceOwnerCtx,
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
		},
		{
			testName: "when instance token for own instance should return unauthZ error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID(),
			},
			inputContext:      instanceOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid id should return not found error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID() + "invalid",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "Errors.Instance.NotFound (COMMA-AE3GS)",
		},
		{
			testName: "when delete succeeds should return deletion date",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID(),
			},
			inputContext:       ctxWithSysAuthZ,
			expectedInstanceID: inst.ID(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2.DeleteInstance(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())
			if tc.expectedErrorMsg == "" {
				require.NotNil(t, res)
				require.NotEmpty(t, res.GetDeletionDate())
			}
		})
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
