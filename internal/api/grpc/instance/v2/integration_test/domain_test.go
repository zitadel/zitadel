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
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func TestAddCustomDomain(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId:   inst.ID(),
				CustomDomain: integration.DomainName(),
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.AddCustomDomainRequest{
				InstanceId:   inst.ID(),
				CustomDomain: integration.DomainName(),
			},
			inputContext:      iamOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.AddCustomDomainRequest{
				InstanceId:   inst.ID(),
				CustomDomain: " ",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (INST-28nlD)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.AddCustomDomainRequest{
				InstanceId:   inst.ID(),
				CustomDomain: " " + integration.DomainName(),
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Cleanup(func() {
				if tc.expectedErrorMsg == "" {
					inst.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{CustomDomain: strings.TrimSpace(tc.inputRequest.CustomDomain)})
				}
			})

			// Test
			res, err := inst.Client.InstanceV2.AddCustomDomain(tc.inputContext, tc.inputRequest)

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
	iamOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)

	customDomain := integration.DomainName()

	_, err := inst.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: inst.ID(), CustomDomain: customDomain})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: inst.ID(), CustomDomain: customDomain})
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId:   inst.ID(),
				CustomDomain: "custom1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				InstanceId:   inst.ID(),
				CustomDomain: "custom1",
			},
			inputContext:      iamOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.RemoveCustomDomainRequest{
				InstanceId:   inst.ID(),
				CustomDomain: " ",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (INST-39nls)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.RemoveCustomDomainRequest{
				InstanceId:   inst.ID(),
				CustomDomain: " " + customDomain,
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2.RemoveCustomDomain(tc.inputContext, tc.inputRequest)

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
	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId:    inst.ID(),
				TrustedDomain: "trusted1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId:    inst.ID(),
				TrustedDomain: "trusted1",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			testName: "when invalid domain should return invalid argument error",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId:    inst.ID(),
				TrustedDomain: " ",
			},
			inputContext:      ctxWithSysAuthZ,
			expectedErrorCode: codes.InvalidArgument,
			expectedErrorMsg:  "Errors.Invalid.Argument (COMMA-Stk21)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.AddTrustedDomainRequest{
				InstanceId:    inst.ID(),
				TrustedDomain: " " + integration.DomainName(),
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Cleanup(func() {
				if tc.expectedErrorMsg == "" {
					inst.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{TrustedDomain: strings.TrimSpace(tc.inputRequest.TrustedDomain)})
				}
			})

			// Test
			res, err := inst.Client.InstanceV2.AddTrustedDomain(tc.inputContext, tc.inputRequest)

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
	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	trustedDomain := integration.DomainName()

	_, err := inst.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: inst.ID(), TrustedDomain: trustedDomain})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: inst.ID(), TrustedDomain: trustedDomain})
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				InstanceId:    inst.ID(),
				TrustedDomain: "trusted1",
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.RemoveTrustedDomainRequest{
				InstanceId:    inst.ID(),
				TrustedDomain: "trusted1",
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			testName: "when valid request should return successful response",
			inputRequest: &instance.RemoveTrustedDomainRequest{
				InstanceId:    inst.ID(),
				TrustedDomain: " " + trustedDomain,
			},
			inputContext: ctxWithSysAuthZ,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2.RemoveTrustedDomain(tc.inputContext, tc.inputRequest)

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
