//go:build integration

package instance_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDeleteInstace(t *testing.T) {
	// Given
	inst := integration.NewInstance(CTXWithSysAuthZ)

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
				InstanceId: " ",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "instance_id must not be empty (instance_id)",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID() + "invalid",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "Instance not found (COMMA-AE3GS)",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.DeleteInstanceRequest{
				InstanceId: inst.ID(),
			},
			inputContext:       CTXWithSysAuthZ,
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
	// Given
	inst := integration.NewInstance(CTXWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)

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
				InstanceName: " ",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceName: " ",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid input should return invalid argument error",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceName: " ",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "instance_name must not be empty (instance_name)",
		},
		{
			testName: "when update succeeds should change instance name",
			inputRequest: &instance.UpdateInstanceRequest{
				InstanceName: "an-updated-name",
			},
			inputContext:    CTXWithSysAuthZ,
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
					retrievedInstance, err := inst.Client.InstanceV2Beta.GetInstance(tc.inputContext, &emptypb.Empty{})
					require.Nil(tt, err)
					assert.Equal(tt, tc.expectedNewName, retrievedInstance.GetInstance().GetName())
				}, retryDuration, tick, "timeout waiting for expected execution result")
			}
		})
	}
}
