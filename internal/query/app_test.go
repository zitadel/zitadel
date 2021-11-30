package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
	"github.com/lib/pq"
)

var (
	expectedAppQuery = regexp.QuoteMeta(`SELECT zitadel.projections.apps.id,` +
		` zitadel.projections.apps.name,` +
		` zitadel.projections.apps.project_id,` +
		` zitadel.projections.apps.creation_date,` +
		` zitadel.projections.apps.change_date,` +
		` zitadel.projections.apps.resource_owner,` +
		` zitadel.projections.apps.state,` +
		` zitadel.projections.apps.sequence,` +
		// api config
		` zitadel.projections.apps_api_configs.app_id,` +
		` zitadel.projections.apps_api_configs.client_id,` +
		` zitadel.projections.apps_api_configs.auth_method,` +
		// oidc config
		` zitadel.projections.apps_oidc_configs.app_id,` +
		` zitadel.projections.apps_oidc_configs.version,` +
		` zitadel.projections.apps_oidc_configs.client_id,` +
		` zitadel.projections.apps_oidc_configs.redirect_uris,` +
		` zitadel.projections.apps_oidc_configs.response_types,` +
		` zitadel.projections.apps_oidc_configs.grant_types,` +
		` zitadel.projections.apps_oidc_configs.application_type,` +
		` zitadel.projections.apps_oidc_configs.auth_method_type,` +
		` zitadel.projections.apps_oidc_configs.post_logout_redirect_uris,` +
		` zitadel.projections.apps_oidc_configs.is_dev_mode,` +
		` zitadel.projections.apps_oidc_configs.access_token_type,` +
		` zitadel.projections.apps_oidc_configs.access_token_role_assertion,` +
		` zitadel.projections.apps_oidc_configs.id_token_role_assertion,` +
		` zitadel.projections.apps_oidc_configs.id_token_userinfo_assertion,` +
		` zitadel.projections.apps_oidc_configs.clock_skew,` +
		` zitadel.projections.apps_oidc_configs.additional_origins` +
		` FROM zitadel.projections.apps` +
		` LEFT JOIN zitadel.projections.apps_api_configs ON zitadel.projections.apps.id = zitadel.projections.apps_api_configs.app_id` +
		` LEFT JOIN zitadel.projections.apps_oidc_configs ON zitadel.projections.apps.id = zitadel.projections.apps_oidc_configs.app_id`)
	expectedAppsQuery = regexp.QuoteMeta(`SELECT zitadel.projections.apps.id,` +
		` zitadel.projections.apps.name,` +
		` zitadel.projections.apps.project_id,` +
		` zitadel.projections.apps.creation_date,` +
		` zitadel.projections.apps.change_date,` +
		` zitadel.projections.apps.resource_owner,` +
		` zitadel.projections.apps.state,` +
		` zitadel.projections.apps.sequence,` +
		// api config
		` zitadel.projections.apps_api_configs.app_id,` +
		` zitadel.projections.apps_api_configs.client_id,` +
		` zitadel.projections.apps_api_configs.auth_method,` +
		// oidc config
		` zitadel.projections.apps_oidc_configs.app_id,` +
		` zitadel.projections.apps_oidc_configs.version,` +
		` zitadel.projections.apps_oidc_configs.client_id,` +
		` zitadel.projections.apps_oidc_configs.redirect_uris,` +
		` zitadel.projections.apps_oidc_configs.response_types,` +
		` zitadel.projections.apps_oidc_configs.grant_types,` +
		` zitadel.projections.apps_oidc_configs.application_type,` +
		` zitadel.projections.apps_oidc_configs.auth_method_type,` +
		` zitadel.projections.apps_oidc_configs.post_logout_redirect_uris,` +
		` zitadel.projections.apps_oidc_configs.is_dev_mode,` +
		` zitadel.projections.apps_oidc_configs.access_token_type,` +
		` zitadel.projections.apps_oidc_configs.access_token_role_assertion,` +
		` zitadel.projections.apps_oidc_configs.id_token_role_assertion,` +
		` zitadel.projections.apps_oidc_configs.id_token_userinfo_assertion,` +
		` zitadel.projections.apps_oidc_configs.clock_skew,` +
		` zitadel.projections.apps_oidc_configs.additional_origins,` +
		` COUNT(*) OVER ()` +
		` FROM zitadel.projections.apps` +
		` LEFT JOIN zitadel.projections.apps_api_configs ON zitadel.projections.apps.id = zitadel.projections.apps_api_configs.app_id` +
		` LEFT JOIN zitadel.projections.apps_oidc_configs ON zitadel.projections.apps.id = zitadel.projections.apps_oidc_configs.app_id`)
	expectedAppIDsQuery = regexp.QuoteMeta(`SELECT zitadel.projections.apps.id` +
		` FROM zitadel.projections.apps` +
		` LEFT JOIN zitadel.projections.apps_api_configs ON zitadel.projections.apps.id = zitadel.projections.apps_api_configs.app_id` +
		` LEFT JOIN zitadel.projections.apps_oidc_configs ON zitadel.projections.apps.id = zitadel.projections.apps_oidc_configs.app_id`)
	expectedProjectIDByAppQuery = regexp.QuoteMeta(`SELECT zitadel.projections.apps.project_id` +
		` FROM zitadel.projections.apps` +
		` LEFT JOIN zitadel.projections.apps_api_configs ON zitadel.projections.apps.id = zitadel.projections.apps_api_configs.app_id` +
		` LEFT JOIN zitadel.projections.apps_oidc_configs ON zitadel.projections.apps.id = zitadel.projections.apps_oidc_configs.app_id`)
	expectedProjectByAppQuery = regexp.QuoteMeta(`SELECT zitadel.projections.projects.id,` +
		` zitadel.projections.projects.creation_date,` +
		` zitadel.projections.projects.change_date,` +
		` zitadel.projections.projects.resource_owner,` +
		` zitadel.projections.projects.state,` +
		` zitadel.projections.projects.sequence,` +
		` zitadel.projections.projects.name,` +
		` zitadel.projections.projects.project_role_assertion,` +
		` zitadel.projections.projects.project_role_check,` +
		` zitadel.projections.projects.has_project_check,` +
		` zitadel.projections.projects.private_labeling_setting` +
		` FROM zitadel.projections.projects` +
		` JOIN zitadel.projections.apps ON zitadel.projections.projects.id = zitadel.projections.apps.project_id` +
		` LEFT JOIN zitadel.projections.apps_api_configs ON zitadel.projections.apps.id = zitadel.projections.apps_api_configs.app_id` +
		` LEFT JOIN zitadel.projections.apps_oidc_configs ON zitadel.projections.apps.id = zitadel.projections.apps_oidc_configs.app_id`)

	appCols = []string{
		"id",
		"name",
		"project_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		// api config
		"app_id",
		"client_id",
		"auth_method",
		// oidc config
		"app_id",
		"version",
		"client_id",
		"redirect_uris",
		"response_types",
		"grant_types",
		"application_type",
		"auth_method_type",
		"post_logout_redirect_uris",
		"is_dev_mode",
		"access_token_type",
		"access_token_role_assertion",
		"id_token_role_assertion",
		"id_token_userinfo_assertion",
		"clock_skew",
		"additional_origins",
	}
	appsCols = append(appCols, "count")
)

func Test_AppsPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareAppsQuery no result",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					nil,
					nil,
				),
			},
			object: &Apps{Apps: []*App{}},
		},
		{
			name:    "prepareAppsQuery only app",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery api app",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							"app-id",
							"api-client-id",
							domain.APIAuthMethodTypePrivateKeyJWT,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						APIConfig: &APIApp{
							ClientID:       "api-client-id",
							AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery oidc app",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           []string{"https://redirect.to/me"},
							ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: []string{"post.logout.ch"},
							IsDevMode:              true,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      []string{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery oidc app AssertIDTokenUserinfo active",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							false,
							false,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           []string{"https://redirect.to/me"},
							ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: []string{"post.logout.ch"},
							IsDevMode:              false,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  false,
							AssertIDTokenRole:      false,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      []string{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery oidc app AssertIDTokenRole active",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							false,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           []string{"https://redirect.to/me"},
							ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: []string{"post.logout.ch"},
							IsDevMode:              true,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      false,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      []string{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery oidc app AssertAccessTokenRole active",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							false,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           []string{"https://redirect.to/me"},
							ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: []string{"post.logout.ch"},
							IsDevMode:              false,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  false,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      []string{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery oidc app IsDevMode active",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Apps: []*App{
					{
						ID:            "app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           []string{"https://redirect.to/me"},
							ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: []string{"post.logout.ch"},
							IsDevMode:              false,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      []string{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery multiple result",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppsQuery,
					appsCols,
					[][]driver.Value{
						{
							"oidc-app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"oidc-app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
						{
							"api-app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							"api-app-id",
							"api-client-id",
							domain.APIAuthMethodTypePrivateKeyJWT,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Apps: []*App{
					{
						ID:            "oidc-app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           []string{"https://redirect.to/me"},
							ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: []string{"post.logout.ch"},
							IsDevMode:              true,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      []string{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
						},
					},
					{
						ID:            "api-app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						APIConfig: &APIApp{
							ClientID:       "api-client-id",
							AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
						},
					},
				},
			},
		},
		{
			name:    "prepareAppsQuery sql err",
			prepare: prepareAppsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedAppsQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_AppPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareAppQuery no result",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*App)(nil),
		},
		{
			name:    "prepareAppQuery found",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedAppQuery,
					appCols,
					[]driver.Value{
						"app-id",
						"app-name",
						"project-id",
						testNow,
						testNow,
						"ro",
						domain.AppStateActive,
						uint64(20211109),
						// api config
						nil,
						nil,
						nil,
						// oidc config
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
			},
		},
		{
			name:    "prepareAppQuery api app",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					appCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							"app-id",
							"api-client-id",
							domain.APIAuthMethodTypePrivateKeyJWT,
							// oidc config
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
							nil,
						},
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
				APIConfig: &APIApp{
					ClientID:       "api-client-id",
					AuthMethodType: domain.APIAuthMethodTypePrivateKeyJWT,
				},
			},
		},
		{
			name:    "prepareAppQuery oidc app",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					appCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           []string{"https://redirect.to/me"},
					ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: []string{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      []string{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
				},
			},
		},
		{
			name:    "prepareAppQuery oidc app IsDevMode inactive",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					appCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           []string{"https://redirect.to/me"},
					ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: []string{"post.logout.ch"},
					IsDevMode:              false,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      []string{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
				},
			},
		},
		{
			name:    "prepareAppQuery oidc app AssertAccessTokenRole inactive",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					appCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							false,
							true,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           []string{"https://redirect.to/me"},
					ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: []string{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  false,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      []string{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
				},
			},
		},
		{
			name:    "prepareAppQuery oidc app AssertIDTokenRole inactive",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					appCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							false,
							true,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           []string{"https://redirect.to/me"},
					ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: []string{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      false,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      []string{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
				},
			},
		},
		{
			name:    "prepareAppQuery oidc app AssertIDTokenUserinfo inactive",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppQuery,
					appCols,
					[][]driver.Value{
						{
							"app-id",
							"app-name",
							"project-id",
							testNow,
							testNow,
							"ro",
							domain.AppStateActive,
							uint64(20211109),
							// api config
							nil,
							nil,
							nil,
							// oidc config
							"app-id",
							domain.OIDCVersionV1,
							"oidc-client-id",
							pq.StringArray{"https://redirect.to/me"},
							pq.Int32Array{int32(domain.OIDCResponseTypeIDTokenToken)},
							pq.Int32Array{int32(domain.OIDCGrantTypeImplicit)},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							pq.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							false,
							1 * time.Second,
							pq.StringArray{"additional.origin"},
						},
					},
				),
			},
			object: &App{
				ID:            "app-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				State:         domain.AppStateActive,
				Sequence:      20211109,
				Name:          "app-name",
				ProjectID:     "project-id",
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           []string{"https://redirect.to/me"},
					ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: []string{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  false,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      []string{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         []string{"https://redirect.to", "additional.origin"},
				},
			},
		},
		{
			name:    "prepareAppQuery sql err",
			prepare: prepareAppQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedAppQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_AppIDsPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareAppIDsQuery no result",
			prepare: prepareAppIDsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppIDsQuery,
					nil,
					nil,
				),
			},
			object: []string{},
		},
		{
			name:    "prepareAppIDsQuery one result",
			prepare: prepareAppIDsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppIDsQuery,
					[]string{"id"},
					[][]driver.Value{
						{
							"app-id",
						},
					},
				),
			},
			object: []string{"app-id"},
		},
		{
			name:    "prepareAppIDsQuery multiple result",
			prepare: prepareAppIDsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppIDsQuery,
					[]string{"id"},
					[][]driver.Value{
						{
							"oidc-app-id",
						},
						{
							"api-app-id",
						},
					},
				),
			},
			object: []string{"oidc-app-id", "api-app-id"},
		},
		{
			name:    "prepareAppIDsQuery sql err",
			prepare: prepareAppIDsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedAppIDsQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_ProjectIDByAppPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareProjectIDByAppQuery no result",
			prepare: prepareProjectIDByAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedProjectIDByAppQuery,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: "",
		},
		{
			name:    "prepareProjectIDByAppQuery one result",
			prepare: prepareProjectIDByAppQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedProjectIDByAppQuery,
					[]string{"project_id"},
					[]driver.Value{"project-id"},
				),
			},
			object: "project-id",
		},
		{
			name:    "prepareProjectIDByAppQuery sql err",
			prepare: prepareProjectIDByAppQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedProjectIDByAppQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_ProjectByAppPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareProjectByAppQuery no result",
			prepare: prepareProjectByAppQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedProjectByAppQuery,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Project)(nil),
		},
		{
			name:    "prepareProjectByAppQuery found",
			prepare: prepareProjectByAppQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedProjectByAppQuery,
					projectCols,
					[]driver.Value{
						"project-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectStateInactive,
						uint64(20211109),
						"project-name",
						true,
						true,
						true,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					},
				),
			},
			object: &Project{
				ID:                     "project-id",
				CreationDate:           testNow,
				ChangeDate:             testNow,
				ResourceOwner:          "ro",
				Sequence:               20211109,
				Name:                   "project-name",
				State:                  domain.ProjectStateInactive,
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
		},
		{
			name:    "prepareProjectByAppQuery found",
			prepare: prepareProjectByAppQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedProjectByAppQuery,
					projectCols,
					[]driver.Value{
						"project-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectStateInactive,
						uint64(20211109),
						"project-name",
						false,
						true,
						true,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					},
				),
			},
			object: &Project{
				ID:                     "project-id",
				CreationDate:           testNow,
				ChangeDate:             testNow,
				ResourceOwner:          "ro",
				Sequence:               20211109,
				Name:                   "project-name",
				State:                  domain.ProjectStateInactive,
				ProjectRoleAssertion:   false,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
		},
		{
			name:    "prepareProjectByAppQuery found",
			prepare: prepareProjectByAppQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedProjectByAppQuery,
					projectCols,
					[]driver.Value{
						"project-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectStateInactive,
						uint64(20211109),
						"project-name",
						true,
						false,
						true,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					},
				),
			},
			object: &Project{
				ID:                     "project-id",
				CreationDate:           testNow,
				ChangeDate:             testNow,
				ResourceOwner:          "ro",
				Sequence:               20211109,
				Name:                   "project-name",
				State:                  domain.ProjectStateInactive,
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       false,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
		},
		{
			name:    "prepareProjectByAppQuery found",
			prepare: prepareProjectByAppQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedProjectByAppQuery,
					projectCols,
					[]driver.Value{
						"project-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectStateInactive,
						uint64(20211109),
						"project-name",
						true,
						true,
						false,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					},
				),
			},
			object: &Project{
				ID:                     "project-id",
				CreationDate:           testNow,
				ChangeDate:             testNow,
				ResourceOwner:          "ro",
				Sequence:               20211109,
				Name:                   "project-name",
				State:                  domain.ProjectStateInactive,
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        false,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
		},
		{
			name:    "prepareProjectByAppQuery sql err",
			prepare: prepareProjectByAppQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedProjectByAppQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
