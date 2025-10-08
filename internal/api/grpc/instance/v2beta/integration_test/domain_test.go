//go:build integration

package instance_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
)

func TestAddCustomDomain(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId: inst.ID(),
				Domain:     integration.DomainName(),
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.AddCustomDomainRequest{
				InstanceId: inst.ID(),
				Domain:     integration.DomainName(),
			},
			inputContext:      iamOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.AddCustomDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " ",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (INST-28nlD)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.AddCustomDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " " + integration.DomainName(),
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Cleanup(func() {
				if tc.expectedErrorMsg == "" {
					inst.Client.InstanceV2Beta.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{Domain: strings.TrimSpace(tc.inputRequest.Domain)})
				}
			})

			// Test
			res, err := inst.Client.InstanceV2Beta.AddCustomDomain(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.GetCreationDate())
			}
		})
	}
}

func TestRemoveCustomDomain(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)
	inst := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)

	customDomain := integration.DomainName()

	_, err := inst.Client.InstanceV2Beta.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: inst.ID(), Domain: customDomain})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: inst.ID(), Domain: customDomain})
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId: inst.ID(),
				Domain:     "custom1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				InstanceId: inst.ID(),
				Domain:     "custom1",
			},
			inputContext:      iamOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " ",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (INST-39nls)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.RemoveCustomDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " " + customDomain,
			},
			inputContext: ctxWithSysAuthZ,
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
				assert.NotEmpty(t, res.GetDeletionDate())
			}
		})
	}
}

func TestAddTrustedDomain(t *testing.T) {
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
		inputContext      context.Context
		inputRequest      *instance.AddTrustedDomainRequest
		expectedErrorMsg  string
		expectedErrorCode codes.Code
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     "trusted1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     "trusted1",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " ",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (COMMA-Stk21)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " " + integration.DomainName(),
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Cleanup(func() {
				if tc.expectedErrorMsg == "" {
					inst.Client.InstanceV2Beta.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{Domain: strings.TrimSpace(tc.inputRequest.Domain)})
				}
			})

			// Test
			res, err := inst.Client.InstanceV2Beta.AddTrustedDomain(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.GetCreationDate())
			}
		})
	}
}

func TestRemoveTrustedDomain(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)
	inst := integration.NewInstance(ctxWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)

	trustedDomain := integration.DomainName()

	_, err := inst.Client.InstanceV2Beta.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: inst.ID(), Domain: trustedDomain})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: inst.ID(), Domain: trustedDomain})
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName          string
		inputContext      context.Context
		inputRequest      *instance.RemoveTrustedDomainRequest
		expectedErrorMsg  string
		expectedErrorCode codes.Code
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.RemoveTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     "trusted1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.RemoveTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     "trusted1",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.RemoveTrustedDomainRequest{
				InstanceId: inst.ID(),
				Domain:     " " + trustedDomain,
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.RemoveTrustedDomain(tc.inputContext, tc.inputRequest)

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
