//go:build integration

package instance_test

import (
	"context"
	"fmt"
	"strings"
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

func TestAddCustomDomain(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)

	relationalInst := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerRelationalCtx := relationalInst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	integration.EnsureInstanceFeature(t, ctxWithSysAuthZ, relationalInst, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		relationalInst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: relationalInst.ID()})
	})

	type instanceAndCtx struct {
		testType string
		instance *integration.Instance
		ctx      context.Context
	}
	testedInstances := []instanceAndCtx{
		{testType: "eventstore", instance: inst, ctx: iamOwnerCtx},
		{testType: "relational", instance: relationalInst, ctx: iamOwnerRelationalCtx},
	}
	for _, instAndCtx := range testedInstances {
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
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: integration.DomainName(),
				},
				inputContext:      context.Background(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				testName: "when unauthZ context should return unauthZ error",
				inputRequest: &instance.AddCustomDomainRequest{
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: integration.DomainName(),
				},
				inputContext:      instAndCtx.ctx,
				expectedErrorCode: codes.PermissionDenied,
				expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
			},
			{
				testName: "when invalid domain should return invalid argument error",
				inputRequest: &instance.AddCustomDomainRequest{
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: " ",
				},
				inputContext:      ctxWithSysAuthZ,
				expectedErrorCode: codes.InvalidArgument,
				expectedErrorMsg:  "Errors.Invalid.Argument",
			},
			{
				testName: "when valid request should return successful response",
				inputRequest: &instance.AddCustomDomainRequest{
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: " " + integration.DomainName(),
				},
				inputContext: ctxWithSysAuthZ,
			},
		}

		for _, tc := range tt {
			t.Run(fmt.Sprintf("%s - %s", instAndCtx.testType, tc.testName), func(t *testing.T) {
				t.Cleanup(func() {
					if tc.expectedErrorMsg == "" {
						instAndCtx.instance.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{CustomDomain: strings.TrimSpace(tc.inputRequest.CustomDomain)})
					}
				})

				// Test
				res, err := instAndCtx.instance.Client.InstanceV2.AddCustomDomain(tc.inputContext, tc.inputRequest)

				// Verify
				assert.Equal(t, tc.expectedErrorCode, status.Code(err))
				assert.Contains(t, status.Convert(err).Message(), tc.expectedErrorMsg)

				if tc.expectedErrorMsg == "" {
					assert.NotNil(t, res)
					assert.NotEmpty(t, res.GetCreationDate())
				}
			})
		}
	}
}

func TestRemoveCustomDomain(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)
	instES := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerCtxES := instES.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)

	instRelational := integration.NewInstance(ctxWithSysAuthZ)
	iamOwnerCtxRelational := instRelational.WithAuthorizationToken(t.Context(), integration.UserTypeIAMOwner)
	integration.EnsureInstanceFeature(t, ctxWithSysAuthZ, instRelational, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})

	customDomainES := integration.DomainName()
	customDomainRelational := integration.DomainName()

	_, err := instES.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: instES.ID(), CustomDomain: customDomainES})
	require.Nil(t, err)
	_, err = instRelational.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: instRelational.ID(), CustomDomain: customDomainRelational})
	require.Nil(t, err)

	t.Cleanup(func() {
		instES.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: instES.ID(), CustomDomain: customDomainES})
		instES.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instES.ID()})
		instRelational.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: instRelational.ID(), CustomDomain: customDomainRelational})
		instRelational.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instRelational.ID()})
	})

	type instanceAndCtx struct {
		testType     string
		instance     *integration.Instance
		ctx          context.Context
		customDomain string
	}
	testedInstances := []instanceAndCtx{
		{testType: "eventstore", instance: instES, ctx: iamOwnerCtxES, customDomain: customDomainES},
		{testType: "relational", instance: instRelational, ctx: iamOwnerCtxRelational, customDomain: customDomainRelational},
	}

	for _, instAndCtx := range testedInstances {
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
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: "custom1",
				},
				inputContext:      context.Background(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				testName: "when unauthZ context should return unauthZ error",
				inputRequest: &instance.RemoveCustomDomainRequest{
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: "custom1",
				},
				inputContext:      instAndCtx.ctx,
				expectedErrorCode: codes.PermissionDenied,
				expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
			},
			{
				testName: "when invalid domain should return invalid argument error",
				inputRequest: &instance.RemoveCustomDomainRequest{
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: " ",
				},
				inputContext:      ctxWithSysAuthZ,
				expectedErrorCode: codes.InvalidArgument,
				expectedErrorMsg:  "Errors.Invalid.Argument",
			},
			{
				testName: "when valid request should return successful response",
				inputRequest: &instance.RemoveCustomDomainRequest{
					InstanceId:   instAndCtx.instance.ID(),
					CustomDomain: " " + instAndCtx.customDomain,
				},
				inputContext: ctxWithSysAuthZ,
			},
		}

		for _, tc := range tt {
			t.Run(fmt.Sprintf("%s - %s", instAndCtx.testType, tc.testName), func(t *testing.T) {
				// Test
				res, err := instAndCtx.instance.Client.InstanceV2.RemoveCustomDomain(tc.inputContext, tc.inputRequest)

				// Verify
				assert.Equal(t, tc.expectedErrorCode, status.Code(err))
				assert.Contains(t, status.Convert(err).Message(), tc.expectedErrorMsg)

				if tc.expectedErrorMsg == "" {
					assert.NotNil(t, res)
					assert.NotEmpty(t, res.GetDeletionDate())
				}
			})
		}
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
