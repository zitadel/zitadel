package convert

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
	filter_pb_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	filter_pb_v2_beta "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
)

func TestAppToPb(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tt := []struct {
		testName      string
		inputQueryApp *query.App
		expectedPbApp *application.Application
	}{
		{
			testName: "full application conversion",
			inputQueryApp: &query.App{
				ID:           "id",
				CreationDate: now,
				ChangeDate:   now,
				State:        domain.AppStateActive,
				Name:         "test-application",
				APIConfig:    &query.APIApp{},
			},
			expectedPbApp: &application.Application{
				ApplicationId: "id",
				CreationDate:  timestamppb.New(now),
				ChangeDate:    timestamppb.New(now),
				State:         application.ApplicationState_APPLICATION_STATE_ACTIVE,
				Name:          "test-application",
				Configuration: &application.Application_ApiConfiguration{
					ApiConfiguration: &application.APIConfiguration{},
				},
			},
		},
		{
			testName:      "nil application",
			inputQueryApp: nil,
			expectedPbApp: &application.Application{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res := AppToPb(tc.inputQueryApp)

			// Then
			assert.Equal(t, tc.expectedPbApp, res)
		})
	}
}

func TestListApplicationsRequestToModel(t *testing.T) {
	t.Parallel()

	validSearchByNameQuery, err := query.NewAppNameSearchQuery(filter.TextMethodPbToQuery(filter_pb_v2_beta.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS), "test")
	require.NoError(t, err)

	validSearchByProjectQuery, err := query.NewAppProjectIDSearchQuery("project1")
	require.NoError(t, err)

	sysDefaults := systemdefaults.SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 150}

	tt := []struct {
		testName string
		req      *application.ListApplicationsRequest

		expectedResponse *query.AppSearchQueries
		expectedError    error
	}{
		{
			testName: "invalid pagination limit",
			req: &application.ListApplicationsRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true, Limit: uint32(sysDefaults.MaxQueryLimit + 1)},
			},
			expectedResponse: nil,
			expectedError:    zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", sysDefaults.MaxQueryLimit+1, sysDefaults.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			testName: "empty request",
			req: &application.ListApplicationsRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
				Filters:    []*application.ApplicationSearchFilter{},
			},
			expectedResponse: &query.AppSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AppColumnID,
				},
				Queries: []query.SearchQuery{},
			},
		},
		{
			testName: "valid request",
			req: &application.ListApplicationsRequest{
				Filters: []*application.ApplicationSearchFilter{
					{
						Filter: &application.ApplicationSearchFilter_ProjectIdFilter{ProjectIdFilter: &application.ProjectIDFilter{ProjectId: "project1"}},
					},
					{
						Filter: &application.ApplicationSearchFilter_NameFilter{NameFilter: &application.ApplicationNameFilter{Name: "test"}},
					},
				},
				SortingColumn: application.ApplicationSorting_APPLICATION_SORT_BY_NAME,
				Pagination:    &filter_pb_v2.PaginationRequest{Asc: true},
			},

			expectedResponse: &query.AppSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AppColumnName,
				},
				Queries: []query.SearchQuery{
					validSearchByProjectQuery,
					validSearchByNameQuery,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			got, err := ListApplicationsRequestToModel(sysDefaults, tc.req)

			// Then
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestApplicationSortingToColumn(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		sorting  application.ApplicationSorting
		expected query.Column
	}{
		{
			name:     "sort by change date",
			sorting:  application.ApplicationSorting_APPLICATION_SORT_BY_CHANGE_DATE,
			expected: query.AppColumnChangeDate,
		},
		{
			name:     "sort by creation date",
			sorting:  application.ApplicationSorting_APPLICATION_SORT_BY_CREATION_DATE,
			expected: query.AppColumnCreationDate,
		},
		{
			name:     "sort by name",
			sorting:  application.ApplicationSorting_APPLICATION_SORT_BY_NAME,
			expected: query.AppColumnName,
		},
		{
			name:     "sort by state",
			sorting:  application.ApplicationSorting_APPLICATION_SORT_BY_STATE,
			expected: query.AppColumnState,
		},
		{
			name:     "sort by ID",
			sorting:  application.ApplicationSorting_APPLICATION_SORT_BY_ID,
			expected: query.AppColumnID,
		},
		{
			name:     "unknown sorting defaults to ID",
			sorting:  application.ApplicationSorting(99),
			expected: query.AppColumnID,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := appSortingToColumn(tc.sorting)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAppStateToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		state    domain.AppState
		expected application.ApplicationState
	}{
		{
			name:     "active state",
			state:    domain.AppStateActive,
			expected: application.ApplicationState_APPLICATION_STATE_ACTIVE,
		},
		{
			name:     "inactive state",
			state:    domain.AppStateInactive,
			expected: application.ApplicationState_APPLICATION_STATE_INACTIVE,
		},
		{
			name:     "removed state",
			state:    domain.AppStateRemoved,
			expected: application.ApplicationState_APPLICATION_STATE_REMOVED,
		},
		{
			name:     "unspecified state",
			state:    domain.AppStateUnspecified,
			expected: application.ApplicationState_APPLICATION_STATE_UNSPECIFIED,
		},
		{
			name:     "unknown state defaults to unspecified",
			state:    domain.AppState(99),
			expected: application.ApplicationState_APPLICATION_STATE_UNSPECIFIED,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := appStateToPb(tc.state)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAppConfigToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		app      *query.App
		expected application.IsApplicationConfiguration
	}{
		{
			name: "OIDC config",
			app: &query.App{
				OIDCConfig: &query.OIDCApp{},
			},
			expected: &application.Application_OidcConfiguration{
				OidcConfiguration: &application.OIDCConfiguration{
					ResponseTypes:      []application.OIDCResponseType{},
					GrantTypes:         []application.OIDCGrantType{},
					ComplianceProblems: []*application.OIDCLocalizedMessage{},
					ClockSkew:          &durationpb.Duration{},
				},
			},
		},
		{
			name: "SAML config",
			app: &query.App{
				SAMLConfig: &query.SAMLApp{},
			},
			expected: &application.Application_SamlConfiguration{
				SamlConfiguration: &application.SAMLConfiguration{},
			},
		},
		{
			name: "API config",
			app: &query.App{
				APIConfig: &query.APIApp{},
			},
			expected: &application.Application_ApiConfiguration{
				ApiConfiguration: &application.APIConfiguration{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := appConfigToPb(tc.app)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLoginVersionToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name          string
		version       *application.LoginVersion
		expectedVer   *domain.LoginVersion
		expectedURI   *string
		expectedError error
	}{
		{
			name:        "nil version",
			version:     nil,
			expectedVer: gu.Ptr(domain.LoginVersionUnspecified),
			expectedURI: gu.Ptr(""),
		},
		{
			name:        "login v1",
			version:     &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
			expectedVer: gu.Ptr(domain.LoginVersion1),
			expectedURI: gu.Ptr(""),
		},
		{
			name:        "login v2 valid URI",
			version:     &application.LoginVersion{Version: &application.LoginVersion_LoginV2{LoginV2: &application.LoginV2{BaseUri: gu.Ptr("https://valid.url")}}},
			expectedVer: gu.Ptr(domain.LoginVersion2),
			expectedURI: gu.Ptr("https://valid.url"),
		},
		{
			name:          "login v2 invalid URI",
			version:       &application.LoginVersion{Version: &application.LoginVersion_LoginV2{LoginV2: &application.LoginV2{BaseUri: gu.Ptr("://invalid")}}},
			expectedVer:   gu.Ptr(domain.LoginVersion2),
			expectedURI:   gu.Ptr("://invalid"),
			expectedError: &url.Error{Op: "parse", URL: "://invalid", Err: errors.New("missing protocol scheme")},
		},
		{
			name:        "unknown version type",
			version:     &application.LoginVersion{},
			expectedVer: gu.Ptr(domain.LoginVersionUnspecified),
			expectedURI: gu.Ptr(""),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			version, uri, err := loginVersionToDomain(tc.version)

			// Then
			assert.Equal(t, tc.expectedVer, version)
			assert.Equal(t, tc.expectedURI, uri)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestLoginVersionToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		version  domain.LoginVersion
		baseURI  *string
		expected *application.LoginVersion
	}{
		{
			name:     "unspecified version",
			version:  domain.LoginVersionUnspecified,
			baseURI:  nil,
			expected: nil,
		},
		{
			name:    "login v1",
			version: domain.LoginVersion1,
			baseURI: nil,
			expected: &application.LoginVersion{
				Version: &application.LoginVersion_LoginV1{
					LoginV1: &application.LoginV1{},
				},
			},
		},
		{
			name:    "login v2",
			version: domain.LoginVersion2,
			baseURI: gu.Ptr("https://example.com"),
			expected: &application.LoginVersion{
				Version: &application.LoginVersion_LoginV2{
					LoginV2: &application.LoginV2{
						BaseUri: gu.Ptr("https://example.com"),
					},
				},
			},
		},
		{
			name:     "unknown version",
			version:  domain.LoginVersion(99),
			baseURI:  nil,
			expected: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := loginVersionToPb(tc.version, tc.baseURI)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAppQueryToModel(t *testing.T) {
	t.Parallel()

	validAppNameSearchQuery, err := query.NewAppNameSearchQuery(query.TextEquals, "test")
	require.NoError(t, err)

	validAppStateSearchQuery, err := query.NewAppStateSearchQuery(domain.AppStateActive)
	require.NoError(t, err)

	tt := []struct {
		name  string
		query *application.ApplicationSearchFilter

		expectedQuery query.SearchQuery
		expectedError error
	}{
		{
			name: "name query",
			query: &application.ApplicationSearchFilter{
				Filter: &application.ApplicationSearchFilter_NameFilter{
					NameFilter: &application.ApplicationNameFilter{
						Name:   "test",
						Method: filter_pb_v2.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
					},
				},
			},
			expectedQuery: validAppNameSearchQuery,
		},
		{
			name: "state query",
			query: &application.ApplicationSearchFilter{
				Filter: &application.ApplicationSearchFilter_StateFilter{
					StateFilter: application.ApplicationState_APPLICATION_STATE_ACTIVE,
				},
			},
			expectedQuery: validAppStateSearchQuery,
		},
		{
			name: "api application only query",
			query: &application.ApplicationSearchFilter{
				Filter: &application.ApplicationSearchFilter_TypeFilter{
					TypeFilter: application.ApplicationType_APPLICATION_TYPE_API,
				},
			},
			expectedQuery: &query.NotNullQuery{
				Column: query.AppAPIConfigColumnAppID,
			},
		},
		{
			name: "oidc application only query",
			query: &application.ApplicationSearchFilter{
				Filter: &application.ApplicationSearchFilter_TypeFilter{
					TypeFilter: application.ApplicationType_APPLICATION_TYPE_OIDC,
				},
			},
			expectedQuery: &query.NotNullQuery{
				Column: query.AppOIDCConfigColumnAppID,
			},
		},
		{
			name: "saml application only query",
			query: &application.ApplicationSearchFilter{
				Filter: &application.ApplicationSearchFilter_TypeFilter{
					TypeFilter: application.ApplicationType_APPLICATION_TYPE_SAML,
				},
			},
			expectedQuery: &query.NotNullQuery{
				Column: query.AppSAMLConfigColumnAppID,
			},
		},
		{
			name:          "invalid query type",
			query:         &application.ApplicationSearchFilter{},
			expectedQuery: nil,
			expectedError: zerrors.ThrowInvalidArgument(nil, "CONV-z2mAGy", "List.Query.Invalid"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result, err := applicationFilterToQuery(tc.query)

			// Then
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedQuery, result)
		})
	}
}

func TestListApplicationKeysRequestToDomain(t *testing.T) {
	t.Parallel()

	resourceOwnerQuery, err := query.NewAuthNKeyResourceOwnerQuery("org1")
	require.NoError(t, err)

	projectIDQuery, err := query.NewAuthNKeyAggregateIDQuery("project1")
	require.NoError(t, err)

	appIDQuery, err := query.NewAuthNKeyObjectIDQuery("app1")
	require.NoError(t, err)

	sysDefaults := systemdefaults.SystemDefaults{DefaultQueryLimit: 100, MaxQueryLimit: 150}

	tt := []struct {
		name string
		req  *application.ListApplicationKeysRequest

		expectedResult *query.AuthNKeySearchQueries
		expectedError  error
	}{
		{
			name: "invalid pagination limit",
			req: &application.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true, Limit: uint32(sysDefaults.MaxQueryLimit + 1)},
			},
			expectedResult: nil,
			expectedError:  zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", sysDefaults.MaxQueryLimit+1, sysDefaults.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name: "empty request",
			req: &application.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
			},
			expectedResult: &query.AuthNKeySearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AuthNKeyColumnID,
				},
				Queries: []query.SearchQuery{},
			},
		},
		{
			name: "only organization id",
			req: &application.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
				Filters: []*application.ApplicationKeySearchFilter{
					{Filter: &application.ApplicationKeySearchFilter_OrganizationIdFilter{
						OrganizationIdFilter: &application.ApplicationKeyOrganizationIDFilter{OrganizationId: "org1"},
					}},
				},
			},
			expectedResult: &query.AuthNKeySearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AuthNKeyColumnID,
				},
				Queries: []query.SearchQuery{
					resourceOwnerQuery,
				},
			},
		},
		{
			name: "only project id",
			req: &application.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
				Filters: []*application.ApplicationKeySearchFilter{
					{Filter: &application.ApplicationKeySearchFilter_ProjectIdFilter{
						ProjectIdFilter: &application.ApplicationKeyProjectIDFilter{ProjectId: "project1"},
					}},
				},
			},
			expectedResult: &query.AuthNKeySearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AuthNKeyColumnID,
				},
				Queries: []query.SearchQuery{
					projectIDQuery,
				},
			},
		},
		{
			name: "only application id",
			req: &application.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
				Filters: []*application.ApplicationKeySearchFilter{
					{Filter: &application.ApplicationKeySearchFilter_ApplicationIdFilter{
						ApplicationIdFilter: &application.ApplicationKeyApplicationIDFilter{ApplicationId: "app1"},
					}},
				},
			},
			expectedResult: &query.AuthNKeySearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AuthNKeyColumnID,
				},
				Queries: []query.SearchQuery{
					appIDQuery,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := ListApplicationKeysRequestToDomain(sysDefaults, tc.req)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestApplicationKeysToPb(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tt := []struct {
		name     string
		input    []*query.AuthNKey
		expected []*application.ApplicationKey
	}{
		{
			name: "multiple keys",
			input: []*query.AuthNKey{
				{
					ID:            "key1",
					AggregateID:   "project1",
					ApplicationID: "app1",
					CreationDate:  now,
					ResourceOwner: "org1",
					Expiration:    now.Add(24 * time.Hour),
					Type:          domain.AuthNKeyTypeJSON,
				},
				{
					ID:            "key2",
					AggregateID:   "project2",
					ApplicationID: "app1",
					CreationDate:  now.Add(-time.Hour),
					ResourceOwner: "org2",
					Expiration:    now.Add(48 * time.Hour),
					Type:          domain.AuthNKeyTypeNONE,
				},
			},
			expected: []*application.ApplicationKey{
				{
					KeyId:          "key1",
					ApplicationId:  "app1",
					ProjectId:      "project1",
					CreationDate:   timestamppb.New(now),
					OrganizationId: "org1",
					ExpirationDate: timestamppb.New(now.Add(24 * time.Hour)),
				},
				{
					KeyId:          "key2",
					ApplicationId:  "app1",
					ProjectId:      "project2",
					CreationDate:   timestamppb.New(now.Add(-time.Hour)),
					OrganizationId: "org2",
					ExpirationDate: timestamppb.New(now.Add(48 * time.Hour)),
				},
			},
		},
		{
			name:     "empty slice",
			input:    []*query.AuthNKey{},
			expected: []*application.ApplicationKey{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []*application.ApplicationKey{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := ApplicationKeysToPb(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
