//go:build integration

package instance_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetInstance(t *testing.T) {
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
		testName           string
		inputContext       context.Context
		expectedInstanceID string
	}{
		{
			testName:           "when unauthN context should instance",
			inputContext:       context.Background(),
			expectedInstanceID: inst.ID(),
		},
		{
			testName:           "when unauthZ context should return instance",
			inputContext:       orgOwnerCtx,
			expectedInstanceID: inst.ID(),
		},
		{
			testName:           "when request succeeds should return matching instance",
			inputContext:       ctxWithSysAuthZ,
			expectedInstanceID: inst.ID(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.GetInstance(tc.inputContext, &emptypb.Empty{})

			// Verify
			require.NoError(t, err)
			assert.Equal(t, tc.expectedInstanceID, res.GetInstance().GetId())
		})
	}
}

func TestListInstances(t *testing.T) {
	t.Parallel()

	// Given
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ctxWithSysAuthZ := integration.WithSystemAuthorization(ctx)

	inst := integration.NewInstance(ctxWithSysAuthZ)
	inst2 := integration.NewInstance(ctxWithSysAuthZ)
	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst2.ID()})
	})

	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)

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
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when valid request with filter should return paginated response",
			inputRequest: &instance.ListInstancesRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.FieldName_FIELD_NAME_CREATION_DATE.Enum(),
				Queries: []*instance.Query{
					{
						Query: &instance.Query_IdQuery{
							IdQuery: &instance.IdsQuery{
								Ids: []string{inst.ID(), inst2.ID()},
							},
						},
					},
				},
			},
			inputContext:      ctxWithSysAuthZ,
			expectedInstances: []string{inst.ID()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.ListInstances(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				require.NotNil(t, res)

				instaceIDs := []string{}
				for _, i := range res.GetInstances() {
					instaceIDs = append(instaceIDs, i.GetId())
				}
				assert.Subset(t, instaceIDs, tc.expectedInstances)
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

	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)
	d1, d2 := "custom."+gofakeit.DomainName(), "custom."+gofakeit.DomainName()

	_, err := inst.Client.InstanceV2Beta.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{Domain: d1})
	require.Nil(t, err)
	_, err = inst.Client.InstanceV2Beta.AddCustomDomain(ctxWithSysAuthZ, &instance.AddCustomDomainRequest{Domain: d2})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{Domain: d1})
		inst.Client.InstanceV2Beta.RemoveCustomDomain(ctxWithSysAuthZ, &instance.RemoveCustomDomainRequest{Domain: d2})
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
	})

	tt := []struct {
		testName        string
		inputRequest    *instance.ListCustomDomainsRequest
		inputContext    context.Context
		expectedDomains []string
	}{
		{
			testName: "when invalid context should paginated empty response",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext: context.Background(),
		},
		{
			testName: "when unauthZ context should return paginated response",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Queries: []*instance.DomainSearchQuery{
					{
						Query: &instance.DomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{Domain: "custom", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    orgOwnerCtx,
			expectedDomains: []string{d1, d2},
		},
		{
			testName: "when valid request with filter should return paginated response",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Queries: []*instance.DomainSearchQuery{
					{
						Query: &instance.DomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{Domain: "custom", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    ctxWithSysAuthZ,
			expectedDomains: []string{d1, d2},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.ListCustomDomains(tc.inputContext, tc.inputRequest)

			// Verify
			require.NotNil(t, res)
			require.NoError(t, err)

			domains := []string{}
			for _, d := range res.GetDomains() {
				domains = append(domains, d.GetDomain())
			}

			assert.Subset(t, domains, tc.expectedDomains)
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

	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)
	d1, d2 := "trusted."+gofakeit.DomainName(), "trusted."+gofakeit.DomainName()

	_, err := inst.Client.InstanceV2Beta.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{Domain: d1})
	require.Nil(t, err)
	_, err = inst.Client.InstanceV2Beta.AddTrustedDomain(ctxWithSysAuthZ, &instance.AddTrustedDomainRequest{Domain: d2})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{Domain: d1})
		inst.Client.InstanceV2Beta.RemoveTrustedDomain(ctxWithSysAuthZ, &instance.RemoveTrustedDomainRequest{Domain: d2})
		inst.Client.InstanceV2Beta.DeleteInstance(ctxWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.ListTrustedDomainsRequest{
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
		},
		{
			testName: "when valid request with filter should return paginated response",
			inputRequest: &instance.ListTrustedDomainsRequest{
				Pagination:    &filter.PaginationRequest{Offset: 0, Limit: 10},
				SortingColumn: instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
				Queries: []*instance.TrustedDomainSearchQuery{
					{
						Query: &instance.TrustedDomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{Domain: "trusted", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS},
						},
					},
				},
			},
			inputContext:    ctxWithSysAuthZ,
			expectedDomains: []string{d1, d2},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.ListTrustedDomains(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				require.NotNil(t, res)

				domains := []string{}
				for _, d := range res.GetTrustedDomain() {
					domains = append(domains, d.GetDomain())
				}

				assert.Subset(t, domains, tc.expectedDomains)
			}
		})
	}
}
