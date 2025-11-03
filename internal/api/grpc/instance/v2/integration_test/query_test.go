//go:build integration

package instance_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestGetInstance(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)
	inst := integration.NewInstance(ctxWithSysAuthZ)
	instanceOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	organizationOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName           string
		inputContext       context.Context
		inputInstanceID    string
		expectedInstanceID string
		expectedErrorMsg   string
		expectedErrorCode  codes.Code
	}{
		{
			testName:          "when unauthN context should return unauthN error",
			inputContext:      context.Background(),
			inputInstanceID:   inst.ID(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName:          "when unauthZ context should return unauthZ error",
			inputContext:      organizationOwnerCtx,
			inputInstanceID:   inst.ID(),
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			testName:           "when request succeeds should return matching instance (systemUser)",
			inputContext:       ctxWithSysAuthZ,
			inputInstanceID:    inst.ID(),
			expectedInstanceID: inst.ID(),
		},
		{
			testName:           "when request succeeds should return matching instance (own context)",
			inputContext:       instanceOwnerCtx,
			expectedInstanceID: inst.ID(),
		},
		{
			testName:          "when instance not found should return not found error",
			inputContext:      ctxWithSysAuthZ,
			inputInstanceID:   "invalid",
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "Errors.IAM.NotFound (QUERY-n0wng)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2.GetInstance(tc.inputContext, &instance.GetInstanceRequest{InstanceId: tc.inputInstanceID})

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedInstanceID, res.GetInstance().GetId())
			}
		})
	}
}

func TestListInstances(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	instances := make([]*integration.Instance, 2)
	inst := integration.NewInstance(ctxWithSysAuthZ)
	inst2 := integration.NewInstance(ctxWithSysAuthZ)
	instances[0], instances[1] = inst, inst2

	t.Cleanup(func() {
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst2.ID()})
	})

	// Sort in descending order
	slices.SortFunc(instances, func(i1, i2 *integration.Instance) int {
		res := i1.Instance.Details.CreationDate.AsTime().Compare(i2.Instance.Details.CreationDate.AsTime())
		if res == 0 {
			return res
		}
		return -res
	})

	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	tt := []struct {
		testName          string
		inputRequest      *instance.ListInstancesRequest
		inputContext      context.Context
		expectedErrorMsg  string
		expectedErrorCode codes.Code
		expectedInstances []string
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.ListInstancesRequest{
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.ListInstancesRequest{
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "Errors.Token.Invalid (AUTH-7fs1e)",
		},
		{
			testName: "when valid request with filter should return paginated response",
			inputRequest: &instance.ListInstancesRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.FieldName_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.Filter{
					{
						Filter: &instance.Filter_InIdsFilter{
							InIdsFilter: &filter.InIDsFilter{
								Ids: []string{inst.ID(), inst2.ID()},
							},
						},
					},
				},
			},
			inputContext:      ctxWithSysAuthZ,
			expectedInstances: []string{inst2.ID(), inst.ID()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2.ListInstances(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				require.NotNil(t, res)

				require.Len(t, res.GetInstances(), len(tc.expectedInstances))

				for i, ins := range res.GetInstances() {
					assert.Equal(t, tc.expectedInstances[i], ins.GetId())
				}
			}
		})
	}
}

func TestListCustomDomains(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)
	inst := integration.NewInstance(ctxWithSysAuthZ)

	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	instanceOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	d1, d2 := "custom."+integration.DomainName(), "custom."+integration.DomainName()

	_, err := inst.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: inst.ID(), CustomDomain: d1})
	require.Nil(t, err)
	_, err = inst.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: inst.ID(), CustomDomain: d2})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: inst.ID(), CustomDomain: d1})
		inst.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: inst.ID(), CustomDomain: d2})
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName          string
		inputRequest      *instance.ListCustomDomainsRequest
		inputContext      context.Context
		expectedErrorMsg  string
		expectedErrorCode codes.Code
		expectedDomains   []string
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.ListCustomDomainsRequest{
				InstanceId: inst.ID(),
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.ListCustomDomainsRequest{
				InstanceId:    inst.ID(),
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.CustomDomainFilter{
					{
						Filter: &instance.CustomDomainFilter_DomainFilter{
							DomainFilter: &instance.DomainFilter{Domain: "custom", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			testName: "when valid request with filter should return paginated response (systemUser)",
			inputRequest: &instance.ListCustomDomainsRequest{
				InstanceId:    inst.ID(),
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.CustomDomainFilter{
					{
						Filter: &instance.CustomDomainFilter_DomainFilter{
							DomainFilter: &instance.DomainFilter{Domain: "custom", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    ctxWithSysAuthZ,
			expectedDomains: []string{d1, d2},
		},
		{
			testName: "when valid request with filter should return paginated response (own context)",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.CustomDomainFilter{
					{
						Filter: &instance.CustomDomainFilter_DomainFilter{
							DomainFilter: &instance.DomainFilter{Domain: "custom", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    instanceOwnerCtx,
			expectedDomains: []string{d1},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputContext, time.Minute)
			require.EventuallyWithT(t, func(collect *assert.CollectT) {
				// Test
				res, err := inst.Client.InstanceV2.ListCustomDomains(tc.inputContext, tc.inputRequest)

				// Verify
				assert.Equal(collect, tc.expectedErrorCode, status.Code(err))
				assert.Equal(collect, tc.expectedErrorMsg, status.Convert(err).Message())

				if tc.expectedErrorMsg == "" {
					domains := []string{}
					for _, d := range res.GetDomains() {
						domains = append(domains, d.GetDomain())
					}

					assert.Subset(collect, domains, tc.expectedDomains)
				}
			}, retryDuration, tick)
		})
	}
}

func TestListTrustedDomains(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)
	inst := integration.NewInstance(ctxWithSysAuthZ)

	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	instanceOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	d1, d2 := "trusted."+integration.DomainName(), "trusted."+integration.DomainName()

	_, err := inst.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: inst.ID(), TrustedDomain: d1})
	require.Nil(t, err)
	_, err = inst.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: inst.ID(), TrustedDomain: d2})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: inst.ID(), TrustedDomain: d1})
		inst.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: inst.ID(), TrustedDomain: d2})
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName          string
		inputRequest      *instance.ListTrustedDomainsRequest
		inputContext      context.Context
		expectedErrorMsg  string
		expectedErrorCode codes.Code
		expectedDomains   []string
	}{
		{
			testName: "when invalid context should return unauthN error",
			inputRequest: &instance.ListTrustedDomainsRequest{
				InstanceId: inst.ID(),
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.ListTrustedDomainsRequest{
				InstanceId: inst.ID(),
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.NotFound,
			expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
		},
		{
			testName: "when valid request with filter should return paginated response (systemUser)",
			inputRequest: &instance.ListTrustedDomainsRequest{
				InstanceId:    inst.ID(),
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.TrustedDomainFilter{
					{
						Filter: &instance.TrustedDomainFilter_DomainFilter{
							DomainFilter: &instance.DomainFilter{Domain: "trusted", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    ctxWithSysAuthZ,
			expectedDomains: []string{d1, d2},
		},
		{
			testName: "when valid request with filter should return paginated response (own context)",
			inputRequest: &instance.ListTrustedDomainsRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.TrustedDomainFilter{
					{
						Filter: &instance.TrustedDomainFilter_DomainFilter{
							DomainFilter: &instance.DomainFilter{Domain: "trusted", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    instanceOwnerCtx,
			expectedDomains: []string{d1},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputContext, time.Minute)
			require.EventuallyWithT(t, func(collect *assert.CollectT) {
				// Test
				res, err := inst.Client.InstanceV2.ListTrustedDomains(tc.inputContext, tc.inputRequest)

				// Verify
				assert.Equal(collect, tc.expectedErrorCode, status.Code(err))
				assert.Equal(collect, tc.expectedErrorMsg, status.Convert(err).Message())

				if tc.expectedErrorMsg == "" {
					require.NotNil(t, res)

					domains := []string{}
					for _, d := range res.GetTrustedDomain() {
						domains = append(domains, d.GetDomain())
					}

					assert.Subset(collect, domains, tc.expectedDomains)
				}
			}, retryDuration, tick)
		})
	}
}
