package instance

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func Test_InstancesToPb(t *testing.T) {
	instances := []*query.Instance{
		{
			ID:   "instance1",
			Name: "Instance One",
			Domains: []*query.InstanceDomain{
				{
					Domain:       "example.com",
					IsPrimary:    true,
					IsGenerated:  false,
					Sequence:     1,
					CreationDate: time.Unix(123, 0),
					ChangeDate:   time.Unix(124, 0),
					InstanceID:   "instance1",
				},
			},
			Sequence:     1,
			CreationDate: time.Unix(123, 0),
			ChangeDate:   time.Unix(124, 0),
		},
	}

	want := []*instance.Instance{
		{
			Id:   "instance1",
			Name: "Instance One",
			Domains: []*instance.Domain{
				{
					Domain:       "example.com",
					Primary:      true,
					Generated:    false,
					InstanceId:   "instance1",
					CreationDate: &timestamppb.Timestamp{Seconds: 123},
				},
			},
			Version:      build.Version(),
			ChangeDate:   &timestamppb.Timestamp{Seconds: 124},
			CreationDate: &timestamppb.Timestamp{Seconds: 123},
		},
	}

	got := InstancesToPb(instances)
	assert.Equal(t, want, got)
}

func Test_ListInstancesRequestToModel(t *testing.T) {
	t.Parallel()

	searchInstanceByID, err := query.NewInstanceIDsListSearchQuery("instance1", "instance2")
	require.Nil(t, err)

	tt := []struct {
		testName       string
		inputRequest   *instance.ListInstancesRequest
		maxQueryLimit  uint64
		expectedResult *query.InstanceSearchQueries
		expectedError  error
	}{
		{
			testName:      "when query limit exceeds max query limit should return invalid argument error",
			maxQueryLimit: 1,
			inputRequest: &instance.ListInstancesRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.FieldName_FIELD_NAME_ID.Enum(),
				Queries:       []*instance.Query{{Query: &instance.Query_IdQuery{IdQuery: &instance.IdsQuery{Ids: []string{"instance1", "instance2"}}}}},
			},
			expectedError: zerrors.ThrowInvalidArgumentf(errors.New("given: 10, allowed: 1"), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			testName: "when valid request should return instance search query model",
			inputRequest: &instance.ListInstancesRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.FieldName_FIELD_NAME_ID.Enum(),
				Queries:       []*instance.Query{{Query: &instance.Query_IdQuery{IdQuery: &instance.IdsQuery{Ids: []string{"instance1", "instance2"}}}}},
			},
			expectedResult: &query.InstanceSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         10,
					Asc:           true,
					SortingColumn: query.InstanceColumnID,
				},
				Queries: []query.SearchQuery{searchInstanceByID},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			sysDefaults := systemdefaults.SystemDefaults{MaxQueryLimit: tc.maxQueryLimit}

			got, err := ListInstancesRequestToModel(tc.inputRequest, sysDefaults)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, got)

		})
	}
}

func Test_fieldNameToInstanceColumn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldName instance.FieldName
		want      query.Column
	}{
		{
			name:      "ID field",
			fieldName: instance.FieldName_FIELD_NAME_ID,
			want:      query.InstanceColumnID,
		},
		{
			name:      "Name field",
			fieldName: instance.FieldName_FIELD_NAME_NAME,
			want:      query.InstanceColumnName,
		},
		{
			name:      "Creation Date field",
			fieldName: instance.FieldName_FIELD_NAME_CREATION_DATE,
			want:      query.InstanceColumnCreationDate,
		},
		{
			name:      "Unknown field",
			fieldName: instance.FieldName(99),
			want:      query.Column{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := fieldNameToInstanceColumn(tt.fieldName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_instanceQueryToModel(t *testing.T) {
	t.Parallel()

	searchInstanceByID, err := query.NewInstanceIDsListSearchQuery("instance1")
	require.Nil(t, err)

	searchInstanceByDomain, err := query.NewInstanceDomainsListSearchQuery("example.com")
	require.Nil(t, err)

	tests := []struct {
		name        string
		searchQuery *instance.Query
		want        query.SearchQuery
		wantErr     bool
	}{
		{
			name: "ID Query",
			searchQuery: &instance.Query{
				Query: &instance.Query_IdQuery{
					IdQuery: &instance.IdsQuery{
						Ids: []string{"instance1"},
					},
				},
			},
			want:    searchInstanceByID,
			wantErr: false,
		},
		{
			name: "Domain Query",
			searchQuery: &instance.Query{
				Query: &instance.Query_DomainQuery{
					DomainQuery: &instance.DomainsQuery{
						Domains: []string{"example.com"},
					},
				},
			},
			want:    searchInstanceByDomain,
			wantErr: false,
		},
		{
			name: "Invalid Query",
			searchQuery: &instance.Query{
				Query: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := instanceQueryToModel(tt.searchQuery)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_ListCustomDomainsRequestToModel(t *testing.T) {
	t.Parallel()

	querySearchRes, err := query.NewInstanceDomainDomainSearchQuery(query.TextEquals, "example.com")
	require.Nil(t, err)

	queryGeneratedRes, err := query.NewInstanceDomainGeneratedSearchQuery(false)
	require.Nil(t, err)

	tests := []struct {
		name           string
		inputRequest   *instance.ListCustomDomainsRequest
		maxQueryLimit  uint64
		expectedResult *query.InstanceDomainSearchQueries
		expectedError  error
	}{
		{
			name: "when query limit exceeds max query limit should return invalid argument error",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN,
				Queries: []*instance.DomainSearchQuery{
					{
						Query: &instance.DomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
								Domain: "example.com",
							},
						},
					},
				},
			},
			maxQueryLimit: 1,
			expectedError: zerrors.ThrowInvalidArgumentf(errors.New("given: 10, allowed: 1"), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name: "when valid request should return domain search query model",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY,
				Queries: []*instance.DomainSearchQuery{
					{
						Query: &instance.DomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS, Domain: "example.com"}},
					},
					{
						Query: &instance.DomainSearchQuery_GeneratedQuery{
							GeneratedQuery: &instance.DomainGeneratedQuery{Generated: false}},
					},
				},
			},
			maxQueryLimit: 100,
			expectedResult: &query.InstanceDomainSearchQueries{
				SearchRequest: query.SearchRequest{Offset: 0, Limit: 10, Asc: true, SortingColumn: query.InstanceDomainIsPrimaryCol},
				Queries: []query.SearchQuery{
					querySearchRes,
					queryGeneratedRes,
				},
			},
			expectedError: nil,
		},
		{
			name: "when invalid query should return error",
			inputRequest: &instance.ListCustomDomainsRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED,
				Queries: []*instance.DomainSearchQuery{
					{
						Query: nil,
					},
				},
			},
			maxQueryLimit:  100,
			expectedResult: nil,
			expectedError:  zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sysDefaults := systemdefaults.SystemDefaults{MaxQueryLimit: tt.maxQueryLimit}

			got, err := ListCustomDomainsRequestToModel(tt.inputRequest, sysDefaults)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, got)
		})
	}
}

func Test_ListTrustedDomainsRequestToModel(t *testing.T) {
	t.Parallel()

	querySearchRes, err := query.NewInstanceTrustedDomainDomainSearchQuery(query.TextEquals, "example.com")
	require.Nil(t, err)

	tests := []struct {
		name           string
		inputRequest   *instance.ListTrustedDomainsRequest
		maxQueryLimit  uint64
		expectedResult *query.InstanceTrustedDomainSearchQueries
		expectedError  error
	}{
		{
			name: "when query limit exceeds max query limit should return invalid argument error",
			inputRequest: &instance.ListTrustedDomainsRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN,
				Queries: []*instance.TrustedDomainSearchQuery{
					{
						Query: &instance.TrustedDomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{
								Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
								Domain: "example.com",
							},
						},
					},
				},
			},
			maxQueryLimit: 1,
			expectedError: zerrors.ThrowInvalidArgumentf(errors.New("given: 10, allowed: 1"), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name: "when valid request should return domain search query model",
			inputRequest: &instance.ListTrustedDomainsRequest{
				Pagination:    &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				SortingColumn: instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE,
				Queries: []*instance.TrustedDomainSearchQuery{
					{
						Query: &instance.TrustedDomainSearchQuery_DomainQuery{
							DomainQuery: &instance.DomainQuery{Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS, Domain: "example.com"}},
					},
				},
			},
			maxQueryLimit: 100,
			expectedResult: &query.InstanceTrustedDomainSearchQueries{
				SearchRequest: query.SearchRequest{Offset: 0, Limit: 10, Asc: true, SortingColumn: query.InstanceTrustedDomainCreationDateCol},
				Queries:       []query.SearchQuery{querySearchRes},
			},
			expectedError: nil,
		},
		{
			name: "when invalid query should return error",
			inputRequest: &instance.ListTrustedDomainsRequest{
				Pagination: &filter.PaginationRequest{Limit: 10, Offset: 0, Asc: true},
				Queries: []*instance.TrustedDomainSearchQuery{
					{
						Query: nil,
					},
				},
			},
			maxQueryLimit:  100,
			expectedResult: nil,
			expectedError:  zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sysDefaults := systemdefaults.SystemDefaults{MaxQueryLimit: tt.maxQueryLimit}

			got, err := ListTrustedDomainsRequestToModel(tt.inputRequest, sysDefaults)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, got)
		})
	}
}
