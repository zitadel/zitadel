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
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
)

func TestGetApplication(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	apiAppName := integration.ApplicationName()
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

	samlAppName := integration.ApplicationName()
	createdSAMLApp, errSAMLAppCreation := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      samlAppName,
		CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
			SamlRequest: &app.CreateSAMLApplicationRequest{
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
				Metadata:     &app.CreateSAMLApplicationRequest_MetadataXml{MetadataXml: samlMetadataGen(integration.URL())},
			},
		},
	})
	require.Nil(t, errSAMLAppCreation)

	oidcAppName := integration.ApplicationName()
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

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *app.GetApplicationRequest
		inputCtx     context.Context

		expectedErrorType       codes.Code
		expectedAppName         string
		expectedAppID           string
		expectedApplicationType string
	}{
		{
			testName: "when unknown app ID should return not found error",
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.GetApplicationRequest{
				Id: integration.ID(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when user has no permission should return membership not found error",
			inputCtx: NoPermissionCtx,
			inputRequest: &app.GetApplicationRequest{
				Id: createdApiApp.GetAppId(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when providing API app ID should return valid API app result",
			inputCtx: projectOwnerCtx,
			inputRequest: &app.GetApplicationRequest{
				Id: createdApiApp.GetAppId(),
			},

			expectedAppName:         apiAppName,
			expectedAppID:           createdApiApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_ApiConfig{}),
		},
		{
			testName: "when providing SAML app ID should return valid SAML app result",
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.GetApplicationRequest{
				Id: createdSAMLApp.GetAppId(),
			},

			expectedAppName:         samlAppName,
			expectedAppID:           createdSAMLApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_SamlConfig{}),
		},
		{
			testName: "when providing OIDC app ID should return valid OIDC app result",
			inputCtx: IAMOwnerCtx,
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

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, 30*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.AppV2Beta.GetApplication(tc.inputCtx, tc.inputRequest)

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
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	t.Parallel()

	createdApiApp, apiAppName := createAPIAppWithName(t, IAMOwnerCtx, instance, p.GetId())

	createdDeactivatedApiApp, deactivatedApiAppName := createAPIAppWithName(t, IAMOwnerCtx, instance, p.GetId())
	deactivateApp(t, createdDeactivatedApiApp, p.GetId())

	_, createdSAMLApp, samlAppName := createSAMLAppWithName(t, integration.URL(), p.GetId())

	createdOIDCApp, oidcAppName := createOIDCAppWithName(t, integration.URL(), p.GetId())

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
		inputRequest *app.ListApplicationsRequest
		inputCtx     context.Context

		expectedOrderedList []appWithName
		expectedOrderedKeys func(keys []appWithName) any
		actualOrderedKeys   func(keys []*app.Application) any
	}{
		{
			testName: "when no apps found should return empty list",
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: "another-id",
			},

			expectedOrderedList: []appWithName{},
			expectedOrderedKeys: func(keys []appWithName) any { return keys },
			actualOrderedKeys:   func(keys []*app.Application) any { return keys },
		},
		{
			testName: "when user has no read permission should return empty set",
			inputCtx: NoPermissionCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedOrderedList: []appWithName{},
			expectedOrderedKeys: func(keys []appWithName) any { return keys },
			actualOrderedKeys:   func(keys []*app.Application) any { return keys },
		},
		{
			testName: "when sorting by name should return apps sorted by name in descending order",
			inputCtx: IAMOwnerCtx,
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
			testName: "when user is project owner should return apps sorted by name in ascending order",
			inputCtx: projectOwnerCtx,
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
			inputCtx: IAMOwnerCtx,
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
			inputCtx: IAMOwnerCtx,
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
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:  p.GetId(),
				Pagination: &filter.PaginationRequest{Asc: true},
				Filters: []*app.ApplicationSearchFilter{
					{Filter: &app.ApplicationSearchFilter_StateFilter{StateFilter: app.AppState_APP_STATE_ACTIVE}},
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
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId:  p.GetId(),
				Pagination: &filter.PaginationRequest{Asc: true},
				Filters: []*app.ApplicationSearchFilter{
					{Filter: &app.ApplicationSearchFilter_OidcAppOnly{}},
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
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, 30*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.AppV2Beta.ListApplications(tc.inputCtx, tc.inputRequest)

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
	reqForAPIAppCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	app1, appAPIConfigChangeErr := instancePermissionV2.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                appName1,
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	app2, appAPIConfigChangeErr := instancePermissionV2.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                appName2,
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	app3, appAPIConfigChangeErr := instancePermissionV2.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                appName3,
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *app.ListApplicationsRequest
		inputCtx     context.Context

		expectedCode   codes.Code
		expectedAppIDs []string
	}{
		{
			testName: "when user has no read permission should return empty set",
			inputCtx: instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeNoPermission),
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{},
		},
		{
			testName: "when projectOwner should return full app list",
			inputCtx: projectOwnerCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedCode:   codes.OK,
			expectedAppIDs: []string{app1.GetAppId(), app2.GetAppId(), app3.GetAppId()},
		},
		{
			testName: "when orgOwner should return full app list",
			inputCtx: instancePermissionV2.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{app1.GetAppId(), app2.GetAppId(), app3.GetAppId()},
		},
		{
			testName: "when iamOwner user should return full app list",
			inputCtx: iamOwnerCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{app1.GetAppId(), app2.GetAppId(), app3.GetAppId()},
		},
		{
			testName: "when other projectOwner user should return empty list",
			inputCtx: otherProjectOwnerCtx,
			inputRequest: &app.ListApplicationsRequest{
				ProjectId: p.GetId(),
			},

			expectedAppIDs: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, 5*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instancePermissionV2.Client.AppV2Beta.ListApplications(tc.inputCtx, tc.inputRequest)

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
	createdAppKey := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp.GetAppId(), time.Now().AddDate(0, 0, 1))

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *app.GetApplicationKeyRequest
		inputCtx     context.Context

		expectedErrorType codes.Code
		expectedAppKeyID  string
	}{
		{
			testName: "when unknown app ID should return not found error",
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.GetApplicationKeyRequest{
				Id: integration.ID(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when user has no permission should return membership not found error",
			inputCtx: NoPermissionCtx,
			inputRequest: &app.GetApplicationKeyRequest{
				Id: createdAppKey.GetId(),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when providing API app ID should return valid API app result",
			inputCtx: projectOwnerCtx,
			inputRequest: &app.GetApplicationKeyRequest{
				Id: createdAppKey.GetId(),
			},

			expectedAppKeyID: createdAppKey.GetId(),
		},
		{
			testName: "when user is OrgOwner should return request key",
			inputCtx: OrgOwnerCtx,
			inputRequest: &app.GetApplicationKeyRequest{
				Id:        createdAppKey.GetId(),
				ProjectId: p.GetId(),
			},

			expectedAppKeyID: createdAppKey.GetId(),
		},
		{
			testName: "when user is IAMOwner should return request key",
			inputCtx: OrgOwnerCtx,
			inputRequest: &app.GetApplicationKeyRequest{
				Id:             createdAppKey.GetId(),
				OrganizationId: instance.DefaultOrg.GetId(),
			},

			expectedAppKeyID: createdAppKey.GetId(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, 30*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.AppV2Beta.GetApplicationKey(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(t, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {

					assert.Equal(t, tc.expectedAppKeyID, res.GetId())
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

	appKey1 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp1.GetAppId(), in2Days)
	appKey2 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp1.GetAppId(), in3Days)
	appKey3 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp1.GetAppId(), tomorrow)
	appKey4 := createAppKey(t, IAMOwnerCtx, instance, p.GetId(), createdApiApp2.GetAppId(), tomorrow)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *app.ListApplicationKeysRequest
		deps         func() (projectID, applicationID, organizationID string)
		inputCtx     context.Context

		expectedErrorType  codes.Code
		expectedAppKeysIDs []string
	}{
		{
			testName: "when sorting by expiration date should return keys sorted by expiration date ascending",
			inputCtx: LoginUserCtx,
			inputRequest: &app.ListApplicationKeysRequest{
				ResourceId:    &app.ListApplicationKeysRequest_ProjectId{ProjectId: p.GetId()},
				Pagination:    &filter.PaginationRequest{Asc: true},
				SortingColumn: app.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_EXPIRATION,
			},
			expectedAppKeysIDs: []string{appKey3.GetId(), appKey4.GetId(), appKey1.GetId(), appKey2.GetId()},
		},
		{
			testName: "when sorting by creation date should return keys sorted by creation date descending",
			inputCtx: IAMOwnerCtx,
			inputRequest: &app.ListApplicationKeysRequest{
				ResourceId:    &app.ListApplicationKeysRequest_ProjectId{ProjectId: p.GetId()},
				SortingColumn: app.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_CREATION_DATE,
			},
			expectedAppKeysIDs: []string{appKey4.GetId(), appKey3.GetId(), appKey2.GetId(), appKey1.GetId()},
		},
		{
			testName: "when filtering by app ID should return keys matching app ID sorted by ID",
			inputCtx: projectOwnerCtx,
			inputRequest: &app.ListApplicationKeysRequest{
				Pagination: &filter.PaginationRequest{Asc: true},
				ResourceId: &app.ListApplicationKeysRequest_ApplicationId{ApplicationId: createdApiApp1.GetAppId()},
			},
			expectedAppKeysIDs: []string{appKey1.GetId(), appKey2.GetId(), appKey3.GetId()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, 5*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instance.Client.AppV2Beta.ListApplicationKeys(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(ttt, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {
					require.Len(ttt, res.GetKeys(), len(tc.expectedAppKeysIDs))

					for i, k := range res.GetKeys() {
						assert.Equal(ttt, tc.expectedAppKeysIDs[i], k.GetId())
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

	appKey1 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp1.GetAppId(), in2Days)
	appKey2 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp1.GetAppId(), in3Days)
	appKey3 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp1.GetAppId(), tomorrow)
	appKey4 := createAppKey(t, iamOwnerCtx, instancePermissionV2, p.GetId(), createdApiApp2.GetAppId(), tomorrow)

	t.Parallel()

	tt := []struct {
		testName     string
		inputRequest *app.ListApplicationKeysRequest
		deps         func() (projectID, applicationID, organizationID string)
		inputCtx     context.Context

		expectedErrorType  codes.Code
		expectedAppKeysIDs []string
	}{
		{
			testName: "when sorting by expiration date should return keys sorted by expiration date ascending",
			inputCtx: loginUserCtx,
			inputRequest: &app.ListApplicationKeysRequest{
				Pagination:    &filter.PaginationRequest{Asc: true},
				SortingColumn: app.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_EXPIRATION,
			},
			expectedAppKeysIDs: []string{appKey3.GetId(), appKey4.GetId(), appKey1.GetId(), appKey2.GetId()},
		},
		{
			testName: "when sorting by creation date should return keys sorted by creation date descending",
			inputCtx: iamOwnerCtx,
			inputRequest: &app.ListApplicationKeysRequest{
				SortingColumn: app.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_CREATION_DATE,
			},
			expectedAppKeysIDs: []string{appKey4.GetId(), appKey3.GetId(), appKey2.GetId(), appKey1.GetId()},
		},
		{
			testName: "when filtering by app ID should return keys matching app ID sorted by ID",
			inputCtx: projectOwnerCtx,
			inputRequest: &app.ListApplicationKeysRequest{
				Pagination: &filter.PaginationRequest{Asc: true},
				ResourceId: &app.ListApplicationKeysRequest_ApplicationId{ApplicationId: createdApiApp1.GetAppId()},
			},
			expectedAppKeysIDs: []string{appKey1.GetId(), appKey2.GetId(), appKey3.GetId()},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// t.Parallel()

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.inputCtx, 5*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				// When
				res, err := instancePermissionV2.Client.AppV2Beta.ListApplicationKeys(tc.inputCtx, tc.inputRequest)

				// Then
				require.Equal(ttt, tc.expectedErrorType, status.Code(err))
				if tc.expectedErrorType == codes.OK {
					require.Len(ttt, res.GetKeys(), len(tc.expectedAppKeysIDs))

					for i, k := range res.GetKeys() {
						assert.Equal(ttt, tc.expectedAppKeysIDs[i], k.GetId())
					}
				}
			}, retryDuration, tick)
		})
	}
}
