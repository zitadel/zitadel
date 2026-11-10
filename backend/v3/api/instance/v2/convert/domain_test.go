package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	filter_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	filter_v2beta "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	instance_v2beta "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func TestListCustomDomainsBetaRequestToV2Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName        string
		inputRequest    *instance_v2beta.ListCustomDomainsRequest
		expectedRequest *instance_v2.ListCustomDomainsRequest
	}{
		{
			testName: "with all fields",
			inputRequest: &instance_v2beta.ListCustomDomainsRequest{
				InstanceId: "instance1",
				Pagination: &filter_v2beta.PaginationRequest{
					Offset: 0,
					Limit:  10,
					Asc:    true,
				},
				SortingColumn: instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
				Queries: []*instance_v2beta.DomainSearchQuery{
					{
						Query: &instance_v2beta.DomainSearchQuery_DomainQuery{
							DomainQuery: &instance_v2beta.DomainQuery{
								Domain: "test.com",
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
							},
						},
					},
					{
						Query: &instance_v2beta.DomainSearchQuery_GeneratedQuery{
							GeneratedQuery: &instance_v2beta.DomainGeneratedQuery{
								Generated: true,
							},
						},
					},
					{
						Query: &instance_v2beta.DomainSearchQuery_PrimaryQuery{
							PrimaryQuery: &instance_v2beta.DomainPrimaryQuery{
								Primary: true,
							},
						},
					},
				},
			},
			expectedRequest: &instance_v2.ListCustomDomainsRequest{
				InstanceId: "instance1",
				Pagination: &filter_v2.PaginationRequest{
					Offset: 0,
					Limit:  10,
					Asc:    true,
				},
				SortingColumn: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
				Filters: []*instance_v2.CustomDomainFilter{
					{
						Filter: &instance_v2.CustomDomainFilter_DomainFilter{
							DomainFilter: &instance_v2.DomainFilter{
								Domain: "test.com",
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
							},
						},
					},
					{
						Filter: &instance_v2.CustomDomainFilter_GeneratedFilter{
							GeneratedFilter: true,
						},
					},
					{
						Filter: &instance_v2.CustomDomainFilter_PrimaryFilter{
							PrimaryFilter: true,
						},
					},
				},
			},
		},
		{
			testName:     "empty request",
			inputRequest: &instance_v2beta.ListCustomDomainsRequest{},
			expectedRequest: &instance_v2.ListCustomDomainsRequest{
				Pagination: &filter_v2.PaginationRequest{},
				Filters:    []*instance_v2.CustomDomainFilter{},
			},
		},
		{
			testName: "with sorting only",
			inputRequest: &instance_v2beta.ListCustomDomainsRequest{
				SortingColumn: instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
			},
			expectedRequest: &instance_v2.ListCustomDomainsRequest{
				Pagination:    &filter_v2.PaginationRequest{},
				SortingColumn: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters:       []*instance_v2.CustomDomainFilter{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			got := ListCustomDomainsBetaRequestToV2Request(tc.inputRequest)
			assert.Equal(t, tc.expectedRequest, got)
		})
	}
}

func Test_listCustomDomainsBetaSortingColToV2Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName         string
		inputFieldName   instance_v2beta.DomainFieldName
		expecteFieldName instance_v2.DomainFieldName
	}{
		{
			testName:         "creation date",
			inputFieldName:   instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
			expecteFieldName: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE,
		},
		{
			testName:         "domain",
			inputFieldName:   instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
			expecteFieldName: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
		},
		{
			testName:         "generated",
			inputFieldName:   instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED,
			expecteFieldName: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED,
		},
		{
			testName:         "primary",
			inputFieldName:   instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY,
			expecteFieldName: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY,
		},
		{
			testName:         "unspecified",
			inputFieldName:   instance_v2beta.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED,
			expecteFieldName: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED,
		},
		{
			testName:         "default",
			inputFieldName:   instance_v2beta.DomainFieldName(99),
			expecteFieldName: instance_v2.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			got := listCustomDomainsBetaSortingColToV2Request(tc.inputFieldName)
			assert.Equal(t, tc.expecteFieldName, got)
		})
	}
}

func TestListTrustedDomainsBetaRequestToV2Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName        string
		inputRequest    *instance_v2beta.ListTrustedDomainsRequest
		expectedRequest *instance_v2.ListTrustedDomainsRequest
	}{
		{
			testName: "with all fields",
			inputRequest: &instance_v2beta.ListTrustedDomainsRequest{
				InstanceId: "instance1",
				Pagination: &filter_v2beta.PaginationRequest{
					Offset: 0,
					Limit:  10,
					Asc:    true,
				},
				SortingColumn: instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN,
				Queries: []*instance_v2beta.TrustedDomainSearchQuery{
					{
						Query: &instance_v2beta.TrustedDomainSearchQuery_DomainQuery{
							DomainQuery: &instance_v2beta.DomainQuery{
								Domain: "test.com",
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
							},
						},
					},
				},
			},
			expectedRequest: &instance_v2.ListTrustedDomainsRequest{
				InstanceId: "instance1",
				Pagination: &filter_v2.PaginationRequest{
					Offset: 0,
					Limit:  10,
					Asc:    true,
				},
				SortingColumn: instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN,
				Filters: []*instance_v2.TrustedDomainFilter{
					{
						Filter: &instance_v2.TrustedDomainFilter_DomainFilter{
							DomainFilter: &instance_v2.DomainFilter{
								Domain: "test.com",
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
							},
						},
					},
				},
			},
		},
		{
			testName:     "empty request",
			inputRequest: &instance_v2beta.ListTrustedDomainsRequest{},
			expectedRequest: &instance_v2.ListTrustedDomainsRequest{
				Pagination: &filter_v2.PaginationRequest{},
				Filters:    []*instance_v2.TrustedDomainFilter{},
			},
		},
		{
			testName: "with sorting only",
			inputRequest: &instance_v2beta.ListTrustedDomainsRequest{
				SortingColumn: instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
			},
			expectedRequest: &instance_v2.ListTrustedDomainsRequest{
				Pagination:    &filter_v2.PaginationRequest{},
				SortingColumn: instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
				Filters:       []*instance_v2.TrustedDomainFilter{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			got := ListTrustedDomainsBetaRequestToV2Request(tc.inputRequest)
			assert.Equal(t, tc.expectedRequest, got)
		})
	}
}

func Test_listTrustedDomainsBetaSortingColToV2Request(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName         string
		inputFieldName   instance_v2beta.TrustedDomainFieldName
		expecteFieldName instance_v2.TrustedDomainFieldName
	}{
		{
			testName:         "creation date",
			inputFieldName:   instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
			expecteFieldName: instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
		},
		{
			testName:         "domain",
			inputFieldName:   instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN,
			expecteFieldName: instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN,
		},
		{
			testName:         "unspecified",
			inputFieldName:   instance_v2beta.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED,
			expecteFieldName: instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED,
		},
		{
			testName:         "default",
			inputFieldName:   instance_v2beta.TrustedDomainFieldName(99),
			expecteFieldName: instance_v2.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			got := listTrustedDomainsBetaSortingColToV2Request(tc.inputFieldName)
			assert.Equal(t, tc.expecteFieldName, got)
		})
	}
}

func TestTrustedDomainInstanceDomainListModelToGRPCBetaResponse(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	tt := []struct {
		name     string
		input    []*domain.InstanceDomain
		expected []*instance_v2beta.TrustedDomain
	}{
		{
			name: "single domain",
			input: []*domain.InstanceDomain{
				{
					InstanceID: "instance1",
					CreatedAt:  testTime,
					Domain:     "test.com",
				},
			},
			expected: []*instance_v2beta.TrustedDomain{
				{
					InstanceId:   "instance1",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test.com",
				},
			},
		},
		{
			name: "multiple domains",
			input: []*domain.InstanceDomain{
				{
					InstanceID: "instance1",
					CreatedAt:  testTime,
					Domain:     "test1.com",
				},
				{
					InstanceID: "instance2",
					CreatedAt:  testTime,
					Domain:     "test2.com",
				},
			},
			expected: []*instance_v2beta.TrustedDomain{
				{
					InstanceId:   "instance1",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test1.com",
				},
				{
					InstanceId:   "instance2",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test2.com",
				},
			},
		},
		{
			name:     "empty list",
			input:    []*domain.InstanceDomain{},
			expected: []*instance_v2beta.TrustedDomain{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := TrustedDomainInstanceDomainListModelToGRPCBetaResponse(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestTrustedDomainInstanceDomainListModelToGRPCResponse(t *testing.T) {
	t.Parallel()

	testTime := time.Now()

	tt := []struct {
		name     string
		input    []*domain.InstanceDomain
		expected []*instance_v2.TrustedDomain
	}{
		{
			name: "single domain",
			input: []*domain.InstanceDomain{
				{
					InstanceID: "instance1",
					CreatedAt:  testTime,
					Domain:     "test.com",
				},
			},
			expected: []*instance_v2.TrustedDomain{
				{
					InstanceId:   "instance1",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test.com",
				},
			},
		},
		{
			name: "multiple domains",
			input: []*domain.InstanceDomain{
				{
					InstanceID: "instance1",
					CreatedAt:  testTime,
					Domain:     "test1.com",
				},
				{
					InstanceID: "instance2",
					CreatedAt:  testTime,
					Domain:     "test2.com",
				},
			},
			expected: []*instance_v2.TrustedDomain{
				{
					InstanceId:   "instance1",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test1.com",
				},
				{
					InstanceId:   "instance2",
					CreationDate: timestamppb.New(testTime),
					Domain:       "test2.com",
				},
			},
		},
		{
			name:     "empty list",
			input:    []*domain.InstanceDomain{},
			expected: []*instance_v2.TrustedDomain{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := TrustedDomainInstanceDomainListModelToGRPCResponse(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
