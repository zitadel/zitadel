//go:build integration

package app_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
)

func TestGetApplication(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	apiAppName := integration.ApplicationName()
	createdApiApp, errAPIAppCreation := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      apiAppName,
		ApplicationType: &application.CreateApplicationRequest_ApiRequest{
			ApiRequest: &application.CreateAPIApplicationRequest{
				AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
		},
	})
	require.Nil(t, errAPIAppCreation)

	samlAppName := integration.ApplicationName()
	createdSAMLApp, errSAMLAppCreation := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      samlAppName,
		ApplicationType: &application.CreateApplicationRequest_SamlRequest{
			SamlRequest: &application.CreateSAMLApplicationRequest{
				LoginVersion: &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
				Metadata:     &application.CreateSAMLApplicationRequest_MetadataXml{MetadataXml: samlMetadataGen(integration.URL())},
			},
		},
	})
	require.Nil(t, errSAMLAppCreation)

	oidcAppName := integration.ApplicationName()
	createdOIDCApp, errOIDCAppCreation := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      oidcAppName,
		ApplicationType: &application.CreateApplicationRequest_OidcRequest{
			OidcRequest: &application.CreateOIDCApplicationRequest{
				RedirectUris:           []string{"http://example.com"},
				ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				AppType:                application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
				AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
				PostLogoutRedirectUris: []string{"http://example.com/home"},
				Version:                application.OIDCVersion_OIDC_VERSION_1_0,
				AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
				BackChannelLogoutUri:   "http://example.com/logout",
				LoginVersion:           &application.LoginVersion{Version: &application.LoginVersion_LoginV2{LoginV2: &application.LoginV2{BaseUri: &baseURI}}},
			},
		},
	})
	require.Nil(t, errOIDCAppCreation)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *application.GetApplicationRequest
		inputCtx     context.Context

		expectedErrorType       codes.Code
		expectedAppName         string
		expectedAppID           string
		expectedApplicationType string
	}{
		{
			testName: "when unknown application ID should return not found error",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.GetApplicationRequest{
				ApplicationId: integration.ID(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when user has no permission should return membership not found error",
			inputCtx: NoPermissionCtx,
			inputRequest: &application.GetApplicationRequest{
				ApplicationId: createdApiApp.GetApplicationId(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when providing API application ID should return valid API application result",
			inputCtx: projectOwnerCtx,
			inputRequest: &application.GetApplicationRequest{
				ApplicationId: createdApiApp.GetApplicationId(),
			},

			expectedAppName:         apiAppName,
			expectedAppID:           createdApiApp.GetApplicationId(),
			expectedApplicationType: fmt.Sprintf("%T", &application.Application_ApiConfig{}),
		},
		{
			testName: "when providing SAML application ID should return valid SAML application result",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.GetApplicationRequest{
				ApplicationId: createdSAMLApp.GetApplicationId(),
			},

			expectedAppName:         samlAppName,
			expectedAppID:           createdSAMLApp.GetApplicationId(),
			expectedApplicationType: fmt.Sprintf("%T", &application.Application_SamlConfig{}),
		},
		{
			testName: "when providing OIDC application ID should return valid OIDC application result",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.GetApplicationRequest{
				ApplicationId: createdOIDCApp.GetApplicationId(),
			},

			expectedAppName:         oidcAppName,
			expectedAppID:           createdOIDCApp.GetApplicationId(),
			expectedApplicationType: fmt.Sprintf("%T", &application.Application_OidcConfig{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.ApplicationV2.GetApplication(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(t, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {

					assert.Equal(t, tc.expectedAppID, res.GetApplication().GetId())
					assert.Equal(t, tc.expectedAppName, res.GetApplication().GetName())
					assert.NotZero(t, res.GetApplication().GetCreationDate())
					assert.NotZero(t, res.GetApplication().GetChangeDate())

					appType := fmt.Sprintf("%T", res.GetApplication().GetConfig())
					assert.Equal(t, tc.expectedApplicationType, appType)
				}
			}, retryDuration, tick)
		})
	}
}

func TestListApplications(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	t.Parallel()

	createdApiApp, apiAppName := createAPIAppWithName(t, IAMOwnerCtx, instance, p.GetId())

	createdDeactivatedApiApp, deactivatedApiAppName := createAPIAppWithName(t, IAMOwnerCtx, instance, p.GetId())
	deactivateApp(t, createdDeactivatedApiApp, p.GetId())

	_, createdSAMLApp, samlAppName := createSAMLAppWithName(t, integration.URL(), p.GetId())

	createdOIDCApp, oidcAppName := createOIDCAppWithName(t, integration.URL(), p.GetId())

	type appWithName struct {
		app  *application.CreateApplicationResponse
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
		if a.app.GetApplicationId() < b.app.GetApplicationId() {
			return -1
		}
		if a.app.GetApplicationId() > b.app.GetApplicationId() {
			return 1
		}
		return 0
	})

	appsSortedByCreationDate := []appWithName{
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
		inputRequest *application.ListApplicationsRequest
		inputCtx     context.Context

		expectedOrderedList []appWithName
		expectedOrderedKeys func(keys []appWithName) any
		actualOrderedKeys   func(keys []*application.Application) any
	}{
		{
			testName: "when no apps found should return empty list",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				Queries: []*application.ApplicationSearchQuery{
					{Query: &application.ApplicationKeySearchQuery_ProjectIdFilter{}},
				},
				ProjectId: "another-id",
			},

			expectedOrderedList: []appWithName{},
			expectedOrderedKeys: func(keys []appWithName) any { return keys },
			actualOrderedKeys:   func(keys []*application.Application) any { return keys },
		},
		{
			testName: "when user has no read permission should return empty set",
			inputCtx: NoPermissionCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedOrderedList: []appWithName{},
			expectedOrderedKeys: func(keys []appWithName) any { return keys },
			actualOrderedKeys:   func(keys []*application.Application) any { return keys },
		},
		{
			testName: "when sorting by name should return apps sorted by name in descending order",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: application.ApplicationSorting_APPLICATION_SORT_BY_NAME,
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
			actualOrderedKeys: func(apps []*application.Application) any {
				names := make([]string, len(apps))
				for i, a := range apps {
					names[i] = a.GetName()
				}

				return names
			},
		},

		{
			testName: "when user is project owner should return apps sorted by name in ascending order",
			inputCtx: projectOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: application.ApplicationSorting_APPLICATION_SORT_BY_NAME,
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
			actualOrderedKeys: func(apps []*application.Application) any {
				names := make([]string, len(apps))
				for i, a := range apps {
					names[i] = a.GetName()
				}

				return names
			},
		},

		{
			testName: "when sorting by id should return apps sorted by id in descending order",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: application.ApplicationSorting_APPLICATION_SORT_BY_ID,
				Pagination:    &filter.PaginationRequest{Asc: true},
			},
			expectedOrderedList: appsSortedByID,
			expectedOrderedKeys: func(apps []appWithName) any {
				ids := make([]string, len(apps))
				for i, a := range apps {
					ids[i] = a.app.GetApplicationId()
				}

				return ids
			},
			actualOrderedKeys: func(apps []*application.Application) any {
				ids := make([]string, len(apps))
				for i, a := range apps {
					ids[i] = a.GetId()
				}

				return ids
			},
		},
		{
			testName: "when sorting by creation date should return apps sorted by creation date in descending order",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId:     p.GetId(),
				SortingColumn: application.ApplicationSorting_APPLICATION_SORT_BY_CREATION_DATE,
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
			actualOrderedKeys: func(apps []*application.Application) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.GetCreationDate().AsTime()
				}

				return creationDates
			},
		},
		{
			testName: "when filtering by active apps should return active apps only",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId:  p.GetId(),
				Pagination: &filter.PaginationRequest{Asc: true},
				Filters: []*application.ApplicationSearchFilter{
					{Filter: &application.ApplicationSearchFilter_StateFilter{StateFilter: application.ApplicationState_APPLICATION_STATE_ACTIVE}},
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
			actualOrderedKeys: func(apps []*application.Application) any {
				creationDates := make([]time.Time, len(apps))
				for i, a := range apps {
					creationDates[i] = a.GetCreationDate().AsTime()
				}

				return creationDates
			},
		},
		{
			testName: "when filtering by application type should return apps of matching type only",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId:  p.GetId(),
				Pagination: &filter.PaginationRequest{Asc: true},
				Filters: []*application.ApplicationSearchFilter{
					{Filter: &application.ApplicationSearchFilter_OidcAppOnly{}},
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
			actualOrderedKeys: func(apps []*application.Application) any {
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
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.ApplicationV2.ListApplications(tc.inputCtx, tc.inputRequest)

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

func TestListApplications_WithPermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)
	p, projectOwnerCtx := getProjectAndProjectContext(t, instancePermissionV2, iamOwnerCtx)
	_, otherProjectOwnerCtx := getProjectAndProjectContext(t, instancePermissionV2, iamOwnerCtx)

	appName1, appName2, appName3 := integration.ApplicationName(), integration.ApplicationName(), integration.ApplicationName()
	reqForAPIAppCreation := &application.CreateApplicationRequest_ApiRequest{
		ApiRequest: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	app1, appAPIConfigChangeErr := instancePermissionV2.Client.ApplicationV2.CreateApplication(iamOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            appName1,
		ApplicationType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	app2, appAPIConfigChangeErr := instancePermissionV2.Client.ApplicationV2.CreateApplication(iamOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            appName2,
		ApplicationType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	app3, appAPIConfigChangeErr := instancePermissionV2.Client.ApplicationV2.CreateApplication(iamOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            appName3,
		ApplicationType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *application.ListApplicationsRequest
		inputCtx     context.Context

		expectedCode   codes.Code
		expectedAppIDs []string
	}{
		{
			testName: "when user has no read permission should return empty set",
			inputCtx: instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeNoPermission),
			inputRequest: &application.ListApplicationsRequest{
				Queries: []*application.ApplicationKeySearchQuery{
					{Filter: &application.ApplicationKeySearchQuery_ApplicationIdFilter{
						ApplicationIdFilter: &application.ApplicationKeyApplicationIDQuery{ApplicationId: createdApiApp1.GetApplicationId()}}},
				},
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{},
		},
		{
			testName: "when projectOwner should return full application list",
			inputCtx: projectOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedCode:   codes.OK,
			expectedAppIDs: []string{app1.GetApplicationId(), app2.GetApplicationId(), app3.GetApplicationId()},
		},
		{
			testName: "when orgOwner should return full application list",
			inputCtx: instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			inputRequest: &application.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{app1.GetApplicationId(), app2.GetApplicationId(), app3.GetApplicationId()},
		},
		{
			testName: "when iamOwner user should return full application list",
			inputCtx: iamOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{app1.GetApplicationId(), app2.GetApplicationId(), app3.GetApplicationId()},
		},
		{
			testName: "when other projectOwner user should return empty list",
			inputCtx: otherProjectOwnerCtx,
			inputRequest: &application.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instancePermissionV2.Client.ApplicationV2.ListApplications(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(ttt, tc.expectedCode, status.Code(err))

				if err == nil {
					require.Len(ttt, res.GetApplications(), len(tc.expectedAppIDs))

					resAppIDs := []string{}
					for _, a := range res.GetApplications() {
						resAppIDs = append(resAppIDs, a.GetId())
					}

					assert.ElementsMatch(ttt, tc.expectedAppIDs, resAppIDs)
				}
			}, retryDuration, tick)
		})
	}
}

func TestGetApplicationKey(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)
	createdApiApp := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	createdAppKey := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp.GetApplicationId(), time.Now().AddDate(0, 0, 1))

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *application.GetApplicationKeyRequest
		inputCtx     context.Context

		expectedErrorType codes.Code
		expectedAppKeyID  string
	}{
		{
			testName: "when unknown application ID should return not found error",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.GetApplicationKeyRequest{
				KeyId: integration.ID(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when user has no permission should return membership not found error",
			inputCtx: NoPermissionCtx,
			inputRequest: &application.GetApplicationKeyRequest{
				KeyId: createdAppKey.GetKeyId(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when providing API application ID should return valid API application result",
			inputCtx: projectOwnerCtx,
			inputRequest: &application.GetApplicationKeyRequest{
				KeyId: createdAppKey.GetKeyId(),
			},

			expectedAppKeyID: createdAppKey.GetKeyId(),
		},
		{
			testName: "when user is OrgOwner should return request key",
			inputCtx: OrgOwnerCtx,
			inputRequest: &application.GetApplicationKeyRequest{
				KeyId: createdAppKey.GetKeyId(),
			},

			expectedAppKeyID: createdAppKey.GetKeyId(),
		},
		{
			testName: "when user is IAMOwner should return request key",
			inputCtx: OrgOwnerCtx,
			inputRequest: &application.GetApplicationKeyRequest{
				KeyId: createdAppKey.GetKeyId(),
			},

			expectedAppKeyID: createdAppKey.GetKeyId(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.ApplicationV2.GetApplicationKey(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(t, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {

					assert.Equal(t, tc.expectedAppKeyID, res.GetKeyId())
					assert.NotEmpty(t, res.GetCreationDate())
					assert.NotEmpty(t, res.GetExpirationDate())
				}
			}, retryDuration, tick)
		})
	}
}

func TestListApplicationKeys(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	createdApiApp1 := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	createdApiApp2 := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	tomorrow := time.Now().AddDate(0, 0, 1)
	in2Days := tomorrow.AddDate(0, 0, 1)
	in3Days := in2Days.AddDate(0, 0, 1)

	appKey1 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp1.GetApplicationId(), in2Days)
	appKey2 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp1.GetApplicationId(), in3Days)
	appKey3 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp1.GetApplicationId(), tomorrow)
	appKey4 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp2.GetApplicationId(), tomorrow)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *application.ListApplicationKeysRequest
		deps         func() (projectID, applicationID, organizationID string)
		inputCtx     context.Context

		expectedErrorType  codes.Code
		expectedAppKeysIDs []string
	}{
		{
			testName: "when sorting by expiration date should return keys sorted by expiration date ascending",
			inputCtx: LoginUserCtx,
			inputRequest: &application.ListApplicationKeysRequest{
				Pagination:    &filter.PaginationRequest{Asc: true},
				SortingColumn: application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_EXPIRATION,
				Queries: []*application.ApplicationKeySearchQuery{
					{Filter: &application.ApplicationKeySearchQuery_ProjectIdFilter{ProjectIdFilter: &application.ApplicationKeyProjectIDQuery{ProjectId: p.GetId()}}},
				},
			},
			expectedAppKeysIDs: []string{appKey3.GetKeyId(), appKey4.GetKeyId(), appKey1.GetKeyId(), appKey2.GetKeyId()},
		},
		{
			testName: "when sorting by creation date should return keys sorted by creation date descending",
			inputCtx: IAMOwnerCtx,
			inputRequest: &application.ListApplicationKeysRequest{
				SortingColumn: application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_CREATION_DATE,
				Queries: []*application.ApplicationKeySearchQuery{
					{Filter: &application.ApplicationKeySearchQuery_ProjectIdFilter{
						ProjectIdFilter: &application.ApplicationKeyProjectIDQuery{ProjectId: p.GetId()}}},
				},
			},
			expectedAppKeysIDs: []string{appKey4.GetKeyId(), appKey3.GetKeyId(), appKey2.GetKeyId(), appKey1.GetKeyId()},
		},
		{
			testName: "when filtering by application ID should return keys matching application ID sorted by ID",
			inputCtx: projectOwnerCtx,
			inputRequest: &application.ListApplicationKeysRequest{
				Pagination: &filter.PaginationRequest{Asc: true},
				Queries: []*application.ApplicationKeySearchQuery{
					{Filter: &application.ApplicationKeySearchQuery_ApplicationIdFilter{
						ApplicationIdFilter: &application.ApplicationKeyApplicationIDQuery{ApplicationId: createdApiApp1.GetApplicationId()}}},
				},
			},
			expectedAppKeysIDs: []string{appKey1.GetKeyId(), appKey2.GetKeyId(), appKey3.GetKeyId()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.ApplicationV2.ListApplicationKeys(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(ttt, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {
					require.Len(ttt, res.GetKeys(), len(tc.expectedAppKeysIDs))

					for i, k := range res.GetKeys() {
						assert.Equal(ttt, tc.expectedAppKeysIDs[i], k.GetKeyId())
					}
				}
			}, retryDuration, tick)
		})
	}
}

func TestListApplicationKeys_WithPermissionV2(t *testing.T) {
	ensureFeaturePermissionV2Enabled(t, instancePermissionV2)
	iamOwnerCtx := instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)
	loginUserCtx := instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeLogin)
	p, projectOwnerCtx := getProjectAndProjectContext(t, instancePermissionV2, iamOwnerCtx)

	createdApiApp1 := createAPIApp(t, iamOwnerCtx, instancePermissionV2, p.GetId())
	createdApiApp2 := createAPIApp(t, iamOwnerCtx, instancePermissionV2, p.GetId())

	tomorrow := time.Now().AddDate(0, 0, 1)
	in2Days := tomorrow.AddDate(0, 0, 1)
	in3Days := in2Days.AddDate(0, 0, 1)

	appKey1 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp1.GetApplicationId(), in2Days)
	appKey2 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp1.GetApplicationId(), in3Days)
	appKey3 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp1.GetApplicationId(), tomorrow)
	appKey4 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp2.GetApplicationId(), tomorrow)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *application.ListApplicationKeysRequest
		deps         func() (projectID, applicationID, organizationID string)
		inputCtx     context.Context

		expectedErrorType  codes.Code
		expectedAppKeysIDs []string
	}{
		{
			testName: "when sorting by expiration date should return keys sorted by expiration date ascending",
			inputCtx: loginUserCtx,
			inputRequest: &application.ListApplicationKeysRequest{
				Pagination:    &filter.PaginationRequest{Asc: true},
				SortingColumn: application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_EXPIRATION,
			},
			expectedAppKeysIDs: []string{appKey3.GetKeyId(), appKey4.GetKeyId(), appKey1.GetKeyId(), appKey2.GetKeyId()},
		},
		{
			testName: "when sorting by creation date should return keys sorted by creation date descending",
			inputCtx: iamOwnerCtx,
			inputRequest: &application.ListApplicationKeysRequest{
				SortingColumn: application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_CREATION_DATE,
			},
			expectedAppKeysIDs: []string{appKey4.GetKeyId(), appKey3.GetKeyId(), appKey2.GetKeyId(), appKey1.GetKeyId()},
		},
		{
			testName: "when filtering by application ID should return keys matching application ID sorted by ID",
			inputCtx: projectOwnerCtx,
			inputRequest: &application.ListApplicationKeysRequest{
				Pagination: &filter.PaginationRequest{Asc: true},
				Queries: []*application.ApplicationKeySearchQuery{
					{Filter: &application.ApplicationKeySearchQuery_ApplicationIdFilter{
						ApplicationIdFilter: &application.ApplicationKeyApplicationIDQuery{ApplicationId: createdApiApp1.GetApplicationId()}}},
				},
			},
			expectedAppKeysIDs: []string{appKey1.GetKeyId(), appKey2.GetKeyId(), appKey3.GetKeyId()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instancePermissionV2.Client.ApplicationV2.ListApplicationKeys(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(ttt, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {
					require.Len(ttt, res.GetKeys(), len(tc.expectedAppKeysIDs))

					for i, k := range res.GetKeys() {
						assert.Equal(ttt, tc.expectedAppKeysIDs[i], k.GetKeyId())
					}
				}
			}, retryDuration, tick)
		})
	}
}
