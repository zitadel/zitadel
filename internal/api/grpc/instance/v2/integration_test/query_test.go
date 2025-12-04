//go:build integration

package instance_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
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
	instES := integration.NewInstance(ctxWithSysAuthZ)
	instanceOwnerCtxES := instES.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	organizationOwnerCtxES := instES.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	// Relational
	instRelational := integration.NewInstance(ctxWithSysAuthZ)
	integration.EnsureInstanceFeature(t, ctxWithSysAuthZ, instRelational, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})
	instanceOwnerCtxRelational := instRelational.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	organizationOwnerCtxRelational := instRelational.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	t.Cleanup(func() {
		instES.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instES.ID()})
		instRelational.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instRelational.ID()})
	})

	type instAndCtx struct {
		testType         string
		inst             *integration.Instance
		instanceOwnerCtx context.Context
		orgOwnerCtx      context.Context
	}

	testedInstanceAndCtxs := []instAndCtx{
		{testType: "eventstore", inst: instES, instanceOwnerCtx: instanceOwnerCtxES, orgOwnerCtx: organizationOwnerCtxES},
		{testType: "relational", inst: instRelational, instanceOwnerCtx: instanceOwnerCtxRelational, orgOwnerCtx: organizationOwnerCtxRelational},
	}

	for _, instWithCtx := range testedInstanceAndCtxs {
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
				inputInstanceID:   instWithCtx.inst.ID(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
				testName:          "when unauthZ context should return unauthZ error",
				inputContext:      instWithCtx.orgOwnerCtx,
				inputInstanceID:   instWithCtx.inst.ID(),
				expectedErrorCode: codes.NotFound,
				expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
			},
			{
				testName:           "when request succeeds should return matching instance (systemUser)",
				inputContext:       ctxWithSysAuthZ,
				inputInstanceID:    instWithCtx.inst.ID(),
				expectedInstanceID: instWithCtx.inst.ID(),
			},
			{
				// TODO(IAM-Marco): Decide if we should set the instance in context.
				testName:           "when request succeeds should return matching instance (own context)",
				inputContext:       instWithCtx.instanceOwnerCtx,
				expectedInstanceID: instWithCtx.inst.ID(),
			},
			{
				// TODO(IAM-Marco): This test won't reach the relational side due to an interceptor
				// automatically changing the instance in context to an empty one.
				// The interceptor will look for the instance matching the instance ID passed here.
				// Since the instance doesn't exist, it will put an empty instance in context.
				// But an empty instance has no feature flags enabled.
				testName:          "when instance not found should return not found error",
				inputContext:      ctxWithSysAuthZ,
				inputInstanceID:   "invalid",
				expectedErrorCode: codes.NotFound,
				expectedErrorMsg:  "Errors.IAM.NotFound (QUERY-n0wng)",
			},
		}

		faultyTestCasesForRelational := []string{
			// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
			"when unauthZ context should return unauthZ error",
			// TODO(IAM-Marco): Decide if we should set the instance in context.
			"when request succeeds should return matching instance (own context)",
		}

		for _, tc := range tt {
			if instWithCtx.testType == "relational" && slices.Contains(faultyTestCasesForRelational, tc.testName) {
				continue
			}
			t.Run(fmt.Sprintf("%s - %s", instWithCtx.testType, tc.testName), func(t *testing.T) {
				// Test
				res, err := instWithCtx.inst.Client.InstanceV2.GetInstance(tc.inputContext, &instance.GetInstanceRequest{InstanceId: tc.inputInstanceID})

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
}

func TestListInstances(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	instances := make([]*integration.Instance, 5)
	inst := integration.NewInstance(ctxWithSysAuthZ)
	inst2 := integration.NewInstance(ctxWithSysAuthZ)
	inst3 := integration.NewInstance(ctxWithSysAuthZ)
	inst4 := integration.NewInstance(ctxWithSysAuthZ)
	inst5 := integration.NewInstance(ctxWithSysAuthZ)
	instances[0], instances[1], instances[2], instances[3], instances[4] = inst, inst2, inst3, inst4, inst5

	t.Cleanup(func() {
		inst.Client.FeatureV2.ResetInstanceFeatures(ctxWithSysAuthZ, &feature.ResetInstanceFeaturesRequest{})
		inst.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		inst2.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst2.ID()})
		inst3.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst3.ID()})
		inst4.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst4.ID()})
		inst5.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst5.ID()})
	})

	orgOwnerCtx := inst.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)

	relTableState := integration.RelationalTablesEnableMatrix()

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
			// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
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
				Pagination:    &filter.PaginationRequest{Offset: 1, Limit: 3},
				SortingColumn: instance.FieldName_FIELD_NAME_CREATION_DATE,
				Filters: []*instance.Filter{
					{
						Filter: &instance.Filter_InIdsFilter{
							InIdsFilter: &filter.InIDsFilter{
								Ids: []string{inst.ID(), inst2.ID(), inst3.ID(), inst4.ID(), inst5.ID()},
							},
						},
					},
					{
						Filter: &instance.Filter_CustomDomainsFilter{
							CustomDomainsFilter: &instance.CustomDomainsFilter{
								Domains: []string{inst4.Domain, inst3.Domain, inst2.Domain, inst.Domain},
							},
						},
					},
				},
			},
			inputContext:      ctxWithSysAuthZ,
			expectedInstances: []string{inst3.ID(), inst2.ID(), inst.ID()},
		},
	}

	for _, stateCase := range relTableState {
		integration.EnsureInstanceFeature(t, ctx, inst, stateCase.FeatureSet, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.Equal(tCollect, stateCase.FeatureSet.GetEnableRelationalTables(), got.EnableRelationalTables.GetEnabled())
		})

		for _, tc := range tt {
			// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
			if tc.testName == "when unauthZ context should return unauthZ error" && stateCase.State == "when relational tables are enabled" {
				continue
			}
			t.Run(fmt.Sprintf("%s - %s", stateCase.State, tc.testName), func(t *testing.T) {
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

}

func TestListCustomDomains(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	// Eventstore objects
	instES := integration.NewInstance(ctxWithSysAuthZ)
	orgOwnerCtxES := instES.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	instanceOwnerCtxES := instES.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	d1ES, d2ES := "custom."+integration.DomainName(), "custom."+integration.DomainName()
	_, err := instES.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: instES.ID(), CustomDomain: d1ES})
	require.Nil(t, err)
	_, err = instES.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: instES.ID(), CustomDomain: d2ES})
	require.Nil(t, err)

	// Relational objects
	instRelational := integration.NewInstance(ctxWithSysAuthZ)
	integration.EnsureInstanceFeature(t, ctx, instRelational, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})
	orgOwnerCtxRelational := instRelational.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	instanceOwnerCtxRelational := instRelational.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	d1Relational, d2Relational := "custom."+integration.DomainName(), "custom."+integration.DomainName()
	_, err = instRelational.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: instRelational.ID(), CustomDomain: d1Relational})
	require.Nil(t, err)
	_, err = instRelational.Client.InstanceV2.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{InstanceId: instRelational.ID(), CustomDomain: d2Relational})
	require.Nil(t, err)

	t.Cleanup(func() {
		instES.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: instES.ID(), CustomDomain: d1ES})
		instES.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: instES.ID(), CustomDomain: d2ES})
		instES.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instES.ID()})

		instRelational.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: instRelational.ID(), CustomDomain: d1Relational})
		instRelational.Client.InstanceV2.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{InstanceId: instRelational.ID(), CustomDomain: d2Relational})
		instRelational.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instRelational.ID()})
	})

	testData := []struct {
		testType        string
		inst            *integration.Instance
		sysOwnerCtx     context.Context
		instOwnerCtx    context.Context
		orgOwnerCtx     context.Context
		expectedDomains []string
	}{
		{testType: "eventstore", inst: instES, sysOwnerCtx: ctxWithSysAuthZ, instOwnerCtx: instanceOwnerCtxES, orgOwnerCtx: orgOwnerCtxES, expectedDomains: []string{d1ES, d2ES}},
		{testType: "relational", inst: instRelational, sysOwnerCtx: ctxWithSysAuthZ, instOwnerCtx: instanceOwnerCtxRelational, orgOwnerCtx: orgOwnerCtxRelational, expectedDomains: []string{d1Relational, d2Relational}},
	}
	for _, td := range testData {

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
					InstanceId: td.inst.ID(),
					Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
				},
				inputContext:      context.Background(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
				testName: "when unauthZ context should return unauthZ error",
				inputRequest: &instance.ListCustomDomainsRequest{
					InstanceId:    td.inst.ID(),
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
				inputContext:      td.orgOwnerCtx,
				expectedErrorCode: codes.NotFound,
				expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
			},
			{
				testName: "when valid request with filter should return paginated response (systemUser)",
				inputRequest: &instance.ListCustomDomainsRequest{
					InstanceId:    td.inst.ID(),
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
				inputContext:    td.sysOwnerCtx,
				expectedDomains: td.expectedDomains,
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
				inputContext:    td.instOwnerCtx,
				expectedDomains: []string{td.expectedDomains[0]},
			},
		}

		for _, tc := range tt {
			if tc.testName == "when unauthZ context should return unauthZ error" && td.testType == "relational" {
				continue
			}
			t.Run(fmt.Sprintf("%s - %s", td.testType, tc.testName), func(t *testing.T) {
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputContext, time.Minute)
				require.EventuallyWithT(t, func(collect *assert.CollectT) {
					// Test
					res, err := td.inst.Client.InstanceV2.ListCustomDomains(tc.inputContext, tc.inputRequest)

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
}

func TestListTrustedDomains(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	// Eventstore objects
	instES := integration.NewInstance(ctxWithSysAuthZ)
	orgOwnerCtxES := instES.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	instanceOwnerCtxES := instES.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	d1ES, d2ES := "trusted."+integration.DomainName(), "trusted."+integration.DomainName()
	_, err := instES.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: instES.ID(), TrustedDomain: d1ES})
	require.Nil(t, err)
	_, err = instES.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: instES.ID(), TrustedDomain: d2ES})
	require.Nil(t, err)

	// Relational objects
	instRelational := integration.NewInstance(ctxWithSysAuthZ)
	integration.EnsureInstanceFeature(t, ctx, instRelational, &feature.SetInstanceFeaturesRequest{EnableRelationalTables: gu.Ptr(true)}, func(tCollect *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
		assert.True(tCollect, got.EnableRelationalTables.GetEnabled())
	})
	orgOwnerCtxRelational := instRelational.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner)
	instanceOwnerCtxRelational := instRelational.WithAuthorizationToken(context.Background(), integration.UserTypeIAMOwner)
	d1Relational, d2Relational := "trusted."+integration.DomainName(), "trusted."+integration.DomainName()
	_, err = instRelational.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: instRelational.ID(), TrustedDomain: d1Relational})
	require.Nil(t, err)
	_, err = instRelational.Client.InstanceV2.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{InstanceId: instRelational.ID(), TrustedDomain: d2Relational})
	require.Nil(t, err)

	t.Cleanup(func() {
		instES.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: instES.ID(), TrustedDomain: d1ES})
		instES.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: instES.ID(), TrustedDomain: d2ES})
		instES.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instES.ID()})

		instRelational.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: instRelational.ID(), TrustedDomain: d1Relational})
		instRelational.Client.InstanceV2.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{InstanceId: instRelational.ID(), TrustedDomain: d2Relational})
		instRelational.Client.InstanceV2.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: instRelational.ID()})
	})

	testData := []struct {
		testType        string
		inst            *integration.Instance
		sysOwnerCtx     context.Context
		instOwnerCtx    context.Context
		orgOwnerCtx     context.Context
		expectedDomains []string
	}{
		{testType: "eventstore", inst: instES, sysOwnerCtx: ctxWithSysAuthZ, instOwnerCtx: instanceOwnerCtxES, orgOwnerCtx: orgOwnerCtxES, expectedDomains: []string{d1ES, d2ES}},
		{testType: "relational", inst: instRelational, sysOwnerCtx: ctxWithSysAuthZ, instOwnerCtx: instanceOwnerCtxRelational, orgOwnerCtx: orgOwnerCtxRelational, expectedDomains: []string{d1Relational, d2Relational}},
	}
	for _, td := range testData {
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
					InstanceId: td.inst.ID(),
					Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
				},
				inputContext:      context.Background(),
				expectedErrorCode: codes.Unauthenticated,
				expectedErrorMsg:  "auth header missing",
			},
			{
				// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
				testName: "when unauthZ context should return unauthZ error",
				inputRequest: &instance.ListTrustedDomainsRequest{
					InstanceId: td.inst.ID(),
					Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
				},
				inputContext:      td.orgOwnerCtx,
				expectedErrorCode: codes.NotFound,
				expectedErrorMsg:  "membership not found (AUTHZ-cdgFk)",
			},
			{
				testName: "when valid request with filter should return paginated response (systemUser)",
				inputRequest: &instance.ListTrustedDomainsRequest{
					InstanceId:    td.inst.ID(),
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
				inputContext:    td.sysOwnerCtx,
				expectedDomains: td.expectedDomains,
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
				inputContext:    td.instOwnerCtx,
				expectedDomains: []string{td.expectedDomains[0]},
			},
		}

		for _, tc := range tt {
			// TODO(IAM-Marco): Fix this test for relational case when permission checks are in place (see https://github.com/zitadel/zitadel/issues/10917)
			if tc.testName == "when unauthZ context should return unauthZ error" && td.testType == "relational" {
				continue
			}

			t.Run(fmt.Sprintf("%s - %s", td.testType, tc.testName), func(t *testing.T) {
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputContext, time.Minute)
				require.EventuallyWithT(t, func(collect *assert.CollectT) {
					// Test
					res, err := td.inst.Client.InstanceV2.ListTrustedDomains(tc.inputContext, tc.inputRequest)

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
}
