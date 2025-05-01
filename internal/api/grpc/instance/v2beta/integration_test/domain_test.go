//go:build integration

package instance_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddCustomDomain(t *testing.T) {
	// Given
	inst := integration.NewInstance(CTXWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)

	t.Cleanup(func() {
		_, err := inst.Client.InstanceV2Beta.DeleteInstance(CTXWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		require.Nil(t, err)
	})

	tt := []struct {
		testName          string
		inputContext      context.Context
		inputRequest      *instance.AddCustomDomainRequest
		expectedErrorMsg  string
		expectedErrorCode codes.Code
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.AddCustomDomainRequest{
				Domain: "custom1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.AddCustomDomainRequest{
				Domain: "custom1",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.AddCustomDomainRequest{
				Domain: " ",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (INST-28nlD)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.AddCustomDomainRequest{
				Domain: " custom-domain.one",
			},
			inputContext: CTXWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.AddCustomDomain(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				assert.NotNil(t, res)
				assert.NotNil(t, res.GetDetails())
				assert.NotEmpty(t, res.GetDetails().GetSequence())
				assert.NotEmpty(t, res.GetDetails().GetCreationDate())
				assert.NotEmpty(t, res.GetDetails().GetResourceOwner())
				assert.Empty(t, res.GetDetails().GetChangeDate())
			}
		})
	}
}

func TestRemoveCustomDomain(t *testing.T) {
	// Given
	inst := integration.NewInstance(CTXWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)

	_, err := inst.Client.InstanceV2Beta.AddCustomDomain(CTXWithSysAuthZ, &instance.AddCustomDomainRequest{Domain: "custom-domain.one"})
	require.Nil(t, err)

	t.Cleanup(func() {
		_, err := inst.Client.InstanceV2Beta.DeleteInstance(CTXWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		require.Nil(t, err)
	})

	tt := []struct {
		testName          string
		inputContext      context.Context
		inputRequest      *instance.RemoveCustomDomainRequest
		expectedErrorMsg  string
		expectedErrorCode codes.Code
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				Domain: "custom1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				Domain: "custom1",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				Domain: " ",
			},
			inputContext:      CTXWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (INST-39nls)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.RemoveCustomDomainRequest{
				Domain: " custom-domain.one",
			},
			inputContext: CTXWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.RemoveCustomDomain(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				assert.NotNil(t, res)
				assert.NotNil(t, res.GetDetails())
				assert.NotEmpty(t, res.GetDetails().GetSequence())
				assert.Empty(t, res.GetDetails().GetCreationDate())
				assert.NotEmpty(t, res.GetDetails().GetResourceOwner())
				assert.NotEmpty(t, res.GetDetails().GetChangeDate())
			}
		})
	}
}
