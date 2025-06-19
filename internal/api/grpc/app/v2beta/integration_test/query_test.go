//go:build integration

package instance_test

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
)

func TestGetApplication(t *testing.T) {
	apiAppName := gofakeit.AppName()
	createdApiApp, errAPIAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: Project.GetId(),
		Name:      apiAppName,
		CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
			ApiRequest: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
		},
	})
	require.Nil(t, errAPIAppCreation)

	samlAppName := gofakeit.AppName()
	createdSAMLApp, errSAMLAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: Project.GetId(),
		Name:      samlAppName,
		CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
			SamlRequest: &app.CreateSAMLApplicationRequest{
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
				Metadata:     &app.CreateSAMLApplicationRequest_MetadataXml{MetadataXml: samlMetadataGen(gofakeit.URL())},
			},
		},
	})
	require.Nil(t, errSAMLAppCreation)

	oidcAppName := gofakeit.AppName()
	createdOIDCApp, errOIDCAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: Project.GetId(),
		Name:      oidcAppName,
		CreationRequestType: &app.CreateApplicationRequest_OidcRequest{
			OidcRequest: &app.CreateOIDCApplicationRequest{
				RedirectUris:           []string{"http://example.com"},
				ResponseTypes:          []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:             []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				AppType:                app.OIDCAppType_OIDC_APP_TYPE_WEB,
				AuthMethodType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
				PostLogoutRedirectUris: []string{"http://example.com/home"},
				Version:                app.OIDCVersion_OIDC_VERSION_1_0,
				AccessTokenType:        app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
				BackChannelLogoutUri:   "http://example.com/logout",
				LoginVersion:           &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: &baseURI}}},
			},
		},
	})
	require.Nil(t, errOIDCAppCreation)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *app.GetApplicationRequest

		expectedErrorType       codes.Code
		expectedAppName         string
		expectedAppID           string
		expectedApplicationType string
	}{
		{
			testName: "when unknown app ID should return not found error",
			inputRequest: &app.GetApplicationRequest{
				Id: gofakeit.Sentence(2),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when providing API app ID should return valid API app result",
			inputRequest: &app.GetApplicationRequest{
				Id: createdApiApp.GetAppId(),
			},

			expectedAppName:         apiAppName,
			expectedAppID:           createdApiApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_ApiConfig{}),
		},
		{
			testName: "when providing SAML app ID should return valid SAML app result",
			inputRequest: &app.GetApplicationRequest{
				Id: createdSAMLApp.GetAppId(),
			},

			expectedAppName:         samlAppName,
			expectedAppID:           createdSAMLApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_SamlConfig{}),
		},
		{
			testName: "when providing OIDC app ID should return valid OIDC app result",
			inputRequest: &app.GetApplicationRequest{
				Id: createdOIDCApp.GetAppId(),
			},

			expectedAppName:         oidcAppName,
			expectedAppID:           createdOIDCApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_OidcConfig{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMOwnerCtx, 30*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.AppV2Beta.GetApplication(IAMOwnerCtx, tc.inputRequest)

				// Then
				require.Equal(t, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {

					assert.Equal(t, tc.expectedAppID, res.GetApp().GetId())
					assert.Equal(t, tc.expectedAppName, res.GetApp().GetName())
					assert.NotZero(t, res.GetApp().GetCreationDate())
					assert.NotZero(t, res.GetApp().GetChangeDate())

					appType := fmt.Sprintf("%T", res.GetApp().GetConfig())
					assert.Equal(t, tc.expectedApplicationType, appType)
				}
			}, retryDuration, tick)
		})
	}
}

func TestListApplications(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), gofakeit.Name(), false, false)

	t.Parallel()

	apiAppName := gofakeit.AppName()
	createdApiApp, errAPIAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      apiAppName,
		CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
			ApiRequest: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
		},
	})
	require.Nil(t, errAPIAppCreation)

	deactivatedApiAppName := gofakeit.AppName()
	createdDeactivatedApiApp, errDeactivatedAPIAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      deactivatedApiAppName,
		CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
			ApiRequest: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
		},
	})
	require.Nil(t, errDeactivatedAPIAppCreation)
	_, deactivateErr := instance.Client.AppV2Beta.DeactivateApplication(IAMOwnerCtx, &app.DeactivateApplicationRequest{
		ProjectId: p.GetId(),
		Id:        createdDeactivatedApiApp.GetAppId(),
	})
	require.Nil(t, deactivateErr)

	samlAppName := gofakeit.AppName()
	createdSAMLApp, errSAMLAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      samlAppName,
		CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
			SamlRequest: &app.CreateSAMLApplicationRequest{
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
				Metadata:     &app.CreateSAMLApplicationRequest_MetadataXml{MetadataXml: samlMetadataGen(gofakeit.URL())},
			},
		},
	})
	require.Nil(t, errSAMLAppCreation)

	oidcAppName := gofakeit.AppName()
	createdOIDCApp, errOIDCAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      oidcAppName,
		CreationRequestType: &app.CreateApplicationRequest_OidcRequest{
			OidcRequest: &app.CreateOIDCApplicationRequest{
				RedirectUris:           []string{"http://example.com"},
				ResponseTypes:          []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:             []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				AppType:                app.OIDCAppType_OIDC_APP_TYPE_WEB,
				AuthMethodType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
				PostLogoutRedirectUris: []string{"http://example.com/home"},
				Version:                app.OIDCVersion_OIDC_VERSION_1_0,
				AccessTokenType:        app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
				BackChannelLogoutUri:   "http://example.com/logout",
				LoginVersion:           &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: &baseURI}}},
			},
		},
	})
	require.Nil(t, errOIDCAppCreation)

	type appWithName struct {
		app  *app.CreateApplicationResponse
		name string
	}

	// Sorting
	appsSortedByName := []appWithName{
		{name: apiAppName, app: createdApiApp},
		{name: deactivatedApiAppName, app: createdDeactivatedApiApp},
		{name: samlAppName, app: createdSAMLApp},
		{name: oidcAppName, app: createdOIDCApp},
	}
	slices.SortFunc(appsSortedByName, func(a, b appWithName) int {
		if a.name < b.name {
			return -1
		}
		if a.name > b.name {
			return 1
		}

		return 0
	})

	appsSortedByID := []appWithName{
		{name: apiAppName, app: createdApiApp},
		{name: deactivatedApiAppName, app: createdDeactivatedApiApp},
		{name: samlAppName, app: createdSAMLApp},
		{name: oidcAppName, app: createdOIDCApp},
	}
	slices.SortFunc(appsSortedByID, func(a, b appWithName) int {
		if a.app.GetAppId() < b.app.GetAppId() {
			return -1
		}
		if a.app.GetAppId() > b.app.GetAppId() {
			return 1
		}

		return 0
	})

	appsSortedByCreationDate :=  []appWithName{
		{name: apiAppName, app: createdApiApp},
		{name: deactivatedApiAppName, app: createdDeactivatedApiApp},
		{name: samlAppName, app: createdSAMLApp},
		{name: oidcAppName, app: createdOIDCApp},
	}
	slices.SortFunc(appsSortedByCreationDate, func(a, b appWithName) int {
		aCreationDate := a.app.GetCreationDate().AsTime()
		bCreationDate := b.app.GetCreationDate().AsTime()
		
		if aCreationDate.Before(bCreationDate) {
			return -1
		}
		if bCreationDate.Before(aCreationDate) {
			return 1
		}

		return 0
	})

	tt := []struct {
		testName     string
		inputRequest *app.ListApplicationsRequest

		expectedOrderedList []appWithName
		expectedOrderedKeys func(keys []appWithName) any
		actualOrderedKeys   func(keys []*app.Application) any
	}{
		{
			testName: "when no apps found should return empty list",
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: "another-id",
			},

			expectedOrderedList: []appWithName{},
			expectedOrderedKeys: func(keys []appWithName) any { return keys },
			actualOrderedKeys:   func(keys []*app.Application) any { return keys },
		},
		{
			testName: "when sorting by name should return apps sorted by name in descending order",
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: app.AppSorting_APP_SORT_BY_NAME,
				Pagination:    &filter.PaginationRequest{Asc: true},
			},

			expectedOrderedList: appsSortedByName,
			expectedOrderedKeys: func(apps []appWithName) any {
				names := make([]string, len(apps))
				for i, a := range apps {
					names[i] = a.name
				}

				return names
			},
			actualOrderedKeys: func(apps []*app.Application) any {
				names := make([]string, len(apps))
				for i, a := range apps {
					names[i] = a.GetName()
				}

				return names
			},
		},
		{
			testName: "when sorting by id should return apps sorted by id in descending order",
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: app.AppSorting_APP_SORT_BY_ID,
				Pagination:    &filter.PaginationRequest{Asc: true},
			},
			expectedOrderedList: appsSortedByID,
			expectedOrderedKeys: func(apps []appWithName) any {
				ids := make([]string, len(apps))
				for i, a := range apps {
					ids[i] = a.app.GetAppId()
				}

				return ids
			},
			actualOrderedKeys: func(apps []*app.Application) any {
				ids := make([]string, len(apps))
				for i, a := range apps {
					ids[i] = a.GetId()
				}

				return ids
			},
		},
		{
			testName: "when sorting by creation date should return apps sorted by creation date in descending order",
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: app.AppSorting_APP_SORT_BY_CREATION_DATE,
				Pagination:    &filter.PaginationRequest{Asc: true},
			},
			expectedOrderedList: appsSortedByCreationDate,
			expectedOrderedKeys: func(apps []appWithName) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.app.GetCreationDate().AsTime()
				}

				return creationDates
			},
			actualOrderedKeys: func(apps []*app.Application) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.GetCreationDate().AsTime()
				}

				return creationDates
			},
		},
		{
			testName: "when filtering by active apps should return active apps only",
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:  p.GetId(),
				Pagination: &filter.PaginationRequest{Asc: true},
				Filters: []*app.ApplicationSearchFilter{
					{ApplicationFilter: &app.ApplicationSearchFilter_StateFilter{StateFilter: app.AppState_APP_STATE_ACTIVE}},
				},
			},
			expectedOrderedList: slices.DeleteFunc(
				slices.Clone(appsSortedByID),
				func(a appWithName) bool { return a.name == deactivatedApiAppName },
			),
			expectedOrderedKeys: func(apps []appWithName) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.app.GetCreationDate().AsTime()
				}

				return creationDates
			},
			actualOrderedKeys: func(apps []*app.Application) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.GetCreationDate().AsTime()
				}

				return creationDates
			},
		},
		{
			testName: "when filtering by app type should return apps of matching type only",
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:  p.GetId(),
				Pagination: &filter.PaginationRequest{Asc: true},
				Filters: []*app.ApplicationSearchFilter{
					{ApplicationFilter: &app.ApplicationSearchFilter_OidcAppOnly{}},
				},
			},
			expectedOrderedList: slices.DeleteFunc(
				slices.Clone(appsSortedByID),
				func(a appWithName) bool { return a.name != oidcAppName },
			),
			expectedOrderedKeys: func(apps []appWithName) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.app.GetCreationDate().AsTime()
				}

				return creationDates
			},
			actualOrderedKeys: func(apps []*app.Application) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.GetCreationDate().AsTime()
				}

				return creationDates
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMOwnerCtx, 30*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.AppV2Beta.ListApplications(IAMOwnerCtx, tc.inputRequest)

				// Then
				require.Equal(ttt, codes.OK, status.Code(err))

				if err == nil {
					assert.Len(ttt, res.GetApplications(), len(tc.expectedOrderedList))
					actualOrderedKeys := tc.actualOrderedKeys(res.GetApplications())
					expectedOrderedKeys := tc.expectedOrderedKeys(tc.expectedOrderedList)
					assert.ElementsMatch(ttt, expectedOrderedKeys, actualOrderedKeys)
				}
			}, retryDuration, tick)
		})
	}
}
