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
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	filter_pb_v2 "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	filter_pb_v2_beta "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
)

func TestAppToPb(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tt := []struct {
		testName      string
		inputQueryApp *query.App
		expectedPbApp *app.Application
	}{
		{
			testName: "full app conversion",
			inputQueryApp: &query.App{
				ID:           "id",
				CreationDate: now,
				ChangeDate:   now,
				State:        domain.AppStateActive,
				Name:         "test-app",
				APIConfig:    &query.APIApp{},
			},
			expectedPbApp: &app.Application{
				Id:           "id",
				CreationDate: timestamppb.New(now),
				ChangeDate:   timestamppb.New(now),
				State:        app.AppState_APP_STATE_ACTIVE,
				Name:         "test-app",
				Config: &app.Application_ApiConfig{
					ApiConfig: &app.APIConfig{},
				},
			},
		},
		{
			testName:      "nil app",
			inputQueryApp: nil,
			expectedPbApp: &app.Application{},
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
		req      *app.ListApplicationsRequest

		expectedResponse *query.AppSearchQueries
		expectedError    error
	}{
		{
			testName: "invalid pagination limit",
			req: &app.ListApplicationsRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true, Limit: uint32(sysDefaults.MaxQueryLimit + 1)},
			},
			expectedResponse: nil,
			expectedError:    zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", sysDefaults.MaxQueryLimit+1, sysDefaults.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			testName: "empty request",
			req: &app.ListApplicationsRequest{
				ProjectId:  "project1",
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
			},
			expectedResponse: &query.AppSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AppColumnID,
				},
				Queries: []query.SearchQuery{
					validSearchByProjectQuery,
				},
			},
		},
		{
			testName: "valid request",
			req: &app.ListApplicationsRequest{
				ProjectId: "project1",
				Filters: []*app.ApplicationSearchFilter{
					{
						Filter: &app.ApplicationSearchFilter_NameFilter{NameFilter: &app.ApplicationNameQuery{Name: "test"}},
					},
				},
				SortingColumn: app.AppSorting_APP_SORT_BY_NAME,
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
					validSearchByNameQuery,
					validSearchByProjectQuery,
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

func TestAppSortingToColumn(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		sorting  app.AppSorting
		expected query.Column
	}{
		{
			name:     "sort by change date",
			sorting:  app.AppSorting_APP_SORT_BY_CHANGE_DATE,
			expected: query.AppColumnChangeDate,
		},
		{
			name:     "sort by creation date",
			sorting:  app.AppSorting_APP_SORT_BY_CREATION_DATE,
			expected: query.AppColumnCreationDate,
		},
		{
			name:     "sort by name",
			sorting:  app.AppSorting_APP_SORT_BY_NAME,
			expected: query.AppColumnName,
		},
		{
			name:     "sort by state",
			sorting:  app.AppSorting_APP_SORT_BY_STATE,
			expected: query.AppColumnState,
		},
		{
			name:     "sort by ID",
			sorting:  app.AppSorting_APP_SORT_BY_ID,
			expected: query.AppColumnID,
		},
		{
			name:     "unknown sorting defaults to ID",
			sorting:  app.AppSorting(99),
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
		expected app.AppState
	}{
		{
			name:     "active state",
			state:    domain.AppStateActive,
			expected: app.AppState_APP_STATE_ACTIVE,
		},
		{
			name:     "inactive state",
			state:    domain.AppStateInactive,
			expected: app.AppState_APP_STATE_INACTIVE,
		},
		{
			name:     "removed state",
			state:    domain.AppStateRemoved,
			expected: app.AppState_APP_STATE_REMOVED,
		},
		{
			name:     "unspecified state",
			state:    domain.AppStateUnspecified,
			expected: app.AppState_APP_STATE_UNSPECIFIED,
		},
		{
			name:     "unknown state defaults to unspecified",
			state:    domain.AppState(99),
			expected: app.AppState_APP_STATE_UNSPECIFIED,
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
		expected app.ApplicationConfig
	}{
		{
			name: "OIDC config",
			app: &query.App{
				OIDCConfig: &query.OIDCApp{},
			},
			expected: &app.Application_OidcConfig{
				OidcConfig: &app.OIDCConfig{
					ResponseTypes:      []app.OIDCResponseType{},
					GrantTypes:         []app.OIDCGrantType{},
					ComplianceProblems: []*app.OIDCLocalizedMessage{},
					ClockSkew:          &durationpb.Duration{},
				},
			},
		},
		{
			name: "SAML config",
			app: &query.App{
				SAMLConfig: &query.SAMLApp{},
			},
			expected: &app.Application_SamlConfig{
				SamlConfig: &app.SAMLConfig{
					Metadata: &app.SAMLConfig_MetadataXml{},
				},
			},
		},
		{
			name: "API config",
			app: &query.App{
				APIConfig: &query.APIApp{},
			},
			expected: &app.Application_ApiConfig{
				ApiConfig: &app.APIConfig{},
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
		version       *app.LoginVersion
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
			version:     &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
			expectedVer: gu.Ptr(domain.LoginVersion1),
			expectedURI: gu.Ptr(""),
		},
		{
			name:        "login v2 valid URI",
			version:     &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: gu.Ptr("https://valid.url")}}},
			expectedVer: gu.Ptr(domain.LoginVersion2),
			expectedURI: gu.Ptr("https://valid.url"),
		},
		{
			name:          "login v2 invalid URI",
			version:       &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: gu.Ptr("://invalid")}}},
			expectedVer:   gu.Ptr(domain.LoginVersion2),
			expectedURI:   gu.Ptr("://invalid"),
			expectedError: &url.Error{Op: "parse", URL: "://invalid", Err: errors.New("missing protocol scheme")},
		},
		{
			name:        "unknown version type",
			version:     &app.LoginVersion{},
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
		expected *app.LoginVersion
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
			expected: &app.LoginVersion{
				Version: &app.LoginVersion_LoginV1{
					LoginV1: &app.LoginV1{},
				},
			},
		},
		{
			name:    "login v2",
			version: domain.LoginVersion2,
			baseURI: gu.Ptr("https://example.com"),
			expected: &app.LoginVersion{
				Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{
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
		query *app.ApplicationSearchFilter

		expectedQuery query.SearchQuery
		expectedError error
	}{
		{
			name: "name query",
			query: &app.ApplicationSearchFilter{
				Filter: &app.ApplicationSearchFilter_NameFilter{
					NameFilter: &app.ApplicationNameQuery{
						Name:   "test",
						Method: filter_pb_v2.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
					},
				},
			},
			expectedQuery: validAppNameSearchQuery,
		},
		{
			name: "state query",
			query: &app.ApplicationSearchFilter{
				Filter: &app.ApplicationSearchFilter_StateFilter{
					StateFilter: app.AppState_APP_STATE_ACTIVE,
				},
			},
			expectedQuery: validAppStateSearchQuery,
		},
		{
			name: "api app only query",
			query: &app.ApplicationSearchFilter{
				Filter: &app.ApplicationSearchFilter_ApiAppOnly{},
			},
			expectedQuery: &query.NotNullQuery{
				Column: query.AppAPIConfigColumnAppID,
			},
		},
		{
			name: "oidc app only query",
			query: &app.ApplicationSearchFilter{
				Filter: &app.ApplicationSearchFilter_OidcAppOnly{},
			},
			expectedQuery: &query.NotNullQuery{
				Column: query.AppOIDCConfigColumnAppID,
			},
		},
		{
			name: "saml app only query",
			query: &app.ApplicationSearchFilter{
				Filter: &app.ApplicationSearchFilter_SamlAppOnly{},
			},
			expectedQuery: &query.NotNullQuery{
				Column: query.AppSAMLConfigColumnAppID,
			},
		},
		{
			name:          "invalid query type",
			query:         &app.ApplicationSearchFilter{},
			expectedQuery: nil,
			expectedError: zerrors.ThrowInvalidArgument(nil, "CONV-z2mAGy", "List.Query.Invalid"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result, err := appQueryToModel(tc.query)

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
		req  *app.ListApplicationKeysRequest

		expectedResult *query.AuthNKeySearchQueries
		expectedError  error
	}{
		{
			name: "invalid pagination limit",
			req: &app.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true, Limit: uint32(sysDefaults.MaxQueryLimit + 1)},
			},
			expectedResult: nil,
			expectedError:  zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", sysDefaults.MaxQueryLimit+1, sysDefaults.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name: "empty request",
			req: &app.ListApplicationKeysRequest{
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
			},
			expectedResult: &query.AuthNKeySearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         100,
					Asc:           true,
					SortingColumn: query.AuthNKeyColumnID,
				},
				Queries: nil,
			},
		},
		{
			name: "only organization id",
			req: &app.ListApplicationKeysRequest{
				ResourceId: &app.ListApplicationKeysRequest_OrganizationId{OrganizationId: "org1"},
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
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
			req: &app.ListApplicationKeysRequest{
				ResourceId: &app.ListApplicationKeysRequest_ProjectId{ProjectId: "project1"},
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
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
			req: &app.ListApplicationKeysRequest{
				ResourceId: &app.ListApplicationKeysRequest_ApplicationId{ApplicationId: "app1"},
				Pagination: &filter_pb_v2.PaginationRequest{Asc: true},
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
		expected []*app.ApplicationKey
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
			expected: []*app.ApplicationKey{
				{
					Id:             "key1",
					ApplicationId:  "app1",
					ProjectId:      "project1",
					CreationDate:   timestamppb.New(now),
					OrganizationId: "org1",
					ExpirationDate: timestamppb.New(now.Add(24 * time.Hour)),
				},
				{
					Id:             "key2",
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
			expected: []*app.ApplicationKey{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: []*app.ApplicationKey{},
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
