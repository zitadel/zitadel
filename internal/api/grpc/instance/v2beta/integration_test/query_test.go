//go:build integration

package instance_test

import (
	"context"
	"testing"

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

func TestListInstances(t *testing.T) {
	// Given
	inst := integration.NewInstance(CTXWithSysAuthZ)
	inst2 := integration.NewInstance(CTXWithSysAuthZ)

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
				SortingColumn: instance.FieldName_FIELD_NAME_CREATION_DATE,
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
			inputContext:      CTXWithSysAuthZ,
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
	inst := integration.NewInstance(CTXWithSysAuthZ)
	orgOwnerCtx := inst.WithAuthorization(context.Background(), integration.UserTypeOrgOwner)
	d1, d2 := "custom-domain.one", "custom-domain.two"

	_, err := inst.Client.InstanceV2Beta.AddCustomDomain(CTXWithSysAuthZ, &instance.AddCustomDomainRequest{Domain: d1})
	require.Nil(t, err)
	_, err = inst.Client.InstanceV2Beta.AddCustomDomain(CTXWithSysAuthZ, &instance.AddCustomDomainRequest{Domain: d2})
	require.Nil(t, err)

	t.Cleanup(func() {
		inst.Client.InstanceV2Beta.DeleteInstance(CTXWithSysAuthZ, &instance.DeleteInstanceRequest{InstanceId: inst.ID()})
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
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      context.Background(),
			expectedErrorCode: codes.Unauthenticated,
			expectedErrorMsg:  "auth header missing",
		},
		{
			testName: "when unauthZ context should return unauthZ error",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination: &filter.PaginationRequest{Offset: 0, Limit: 10},
			},
			inputContext:      orgOwnerCtx,
			expectedErrorCode: codes.PermissionDenied,
			expectedErrorMsg:  "No matching permissions found (AUTH-5mWD2)",
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
			inputContext:    CTXWithSysAuthZ,
			expectedDomains: []string{d1, d2},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Test
			res, err := inst.Client.InstanceV2Beta.ListCustomDomains(tc.inputContext, tc.inputRequest)

			// Verify
			assert.Equal(t, tc.expectedErrorCode, status.Code(err))
			assert.Equal(t, tc.expectedErrorMsg, status.Convert(err).Message())

			if tc.expectedErrorMsg == "" {
				require.NotNil(t, res)

				domains := []string{}
				for _, d := range res.GetResult() {
					domains = append(domains, d.GetDomain())
				}

				assert.Subset(t, domains, tc.expectedDomains)
			}
		})
	}
}
