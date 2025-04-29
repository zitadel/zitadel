//go:build integration

package instance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/integration"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetInstance(t *testing.T) {
	inst := integration.NewInstance(CTXWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)
	tt := []struct {
		testName           string
		inputContext       context.Context
		expectedErrorMsg   string
		expectedErrorCode  codes.Code
		expectedInstanceID string
	}{
		{
			testName:          "when invalid context should return unauthN error",
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName:          "when unauthZ context should return unauthZ error",
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName:           "when request succeeds should return matching instance",
			inputContext:       CTXWithSysAuthZ,
			expectedInstanceID: inst.ID(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.GetInstance(tc.inputContext, &emptypb.Empty{})

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				assert.Equal(t, tc.expectedInstanceID, res.GetInstance().GetId())
			}
		})
	}
}
