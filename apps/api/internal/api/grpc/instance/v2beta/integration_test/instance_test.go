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
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func TestDeleteInstace(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId: " ",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID() + "invalid",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "Instance not found (COMMA-AE3GS)",
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
			res, err := inst.Client.InstanceV2Beta.DeleteInstance(tc.inputContext, tc.inputRequest)

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

func TestUpdateInstace(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID(),
				InstanceName: " ",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when update succeeds should change instance name",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceId:   inst.ID(),
				InstanceName: "an-updated-name",
			},
			inputContext:    ctxWithSysAuthZ,
			expectedNewName: "an-updated-name",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.UpdateInstance(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())
			if tc.expectedErrorMsg == "" {

				require.NotNil(t, res)
				assert.NotEmpty(t, res.GetChangeDate())

				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputContext, 20*time.Second)
				require.EventuallyWithT(t, func(tt *assert.CollectT) {
					retrievedInstance, err := inst.Client.InstanceV2Beta.GetInstance(tc.inputContext, &instance.GetInstanceRequest{InstanceId: inst.ID()})
					require.Nil(tt, err)
					assert.Equal(tt, tc.expectedNewName, retrievedInstance.GetInstance().GetName())
				}, retryDuration, tick, "timeout waiting for expected execution result")
			}
		})
	}
}
