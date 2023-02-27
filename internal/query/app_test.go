package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	expectedAppQuery = regexp.QuoteMeta(`SELECT projections.apps4.id,` +
		` projections.apps4.name,` +
		` projections.apps4.project_id,` +
		` projections.apps4.creation_date,` +
		` projections.apps4.change_date,` +
		` projections.apps4.resource_owner,` +
		` projections.apps4.state,` +
		` projections.apps4.sequence,` +
		// api config
		` projections.apps4_api_configs.app_id,` +
		` projections.apps4_api_configs.client_id,` +
		` projections.apps4_api_configs.auth_method,` +
		// oidc config
		` projections.apps4_oidc_configs.app_id,` +
		` projections.apps4_oidc_configs.version,` +
		` projections.apps4_oidc_configs.client_id,` +
		` projections.apps4_oidc_configs.redirect_uris,` +
		` projections.apps4_oidc_configs.response_types,` +
		` projections.apps4_oidc_configs.grant_types,` +
		` projections.apps4_oidc_configs.application_type,` +
		` projections.apps4_oidc_configs.auth_method_type,` +
		` projections.apps4_oidc_configs.post_logout_redirect_uris,` +
		` projections.apps4_oidc_configs.is_dev_mode,` +
		` projections.apps4_oidc_configs.access_token_type,` +
		` projections.apps4_oidc_configs.access_token_role_assertion,` +
		` projections.apps4_oidc_configs.id_token_role_assertion,` +
		` projections.apps4_oidc_configs.id_token_userinfo_assertion,` +
		` projections.apps4_oidc_configs.clock_skew,` +
		` projections.apps4_oidc_configs.additional_origins,` +
		//saml config
		` projections.apps4_saml_configs.app_id,` +
		` projections.apps4_saml_configs.entity_id,` +
		` projections.apps4_saml_configs.metadata,` +
		` projections.apps4_saml_configs.metadata_url` +
		` FROM projections.apps4` +
		` LEFT JOIN projections.apps4_api_configs ON projections.apps4.id = projections.apps4_api_configs.app_id AND projections.apps4.instance_id = projections.apps4_api_configs.instance_id` +
		` LEFT JOIN projections.apps4_oidc_configs ON projections.apps4.id = projections.apps4_oidc_configs.app_id AND projections.apps4.instance_id = projections.apps4_oidc_configs.instance_id` +
		` LEFT JOIN projections.apps4_saml_configs ON projections.apps4.id = projections.apps4_saml_configs.app_id AND projections.apps4.instance_id = projections.apps4_saml_configs.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedAppsQuery = regexp.QuoteMeta(`SELECT projections.apps4.id,` +
		` projections.apps4.name,` +
		` projections.apps4.project_id,` +
		` projections.apps4.creation_date,` +
		` projections.apps4.change_date,` +
		` projections.apps4.resource_owner,` +
		` projections.apps4.state,` +
		` projections.apps4.sequence,` +
		// api config
		` projections.apps4_api_configs.app_id,` +
		` projections.apps4_api_configs.client_id,` +
		` projections.apps4_api_configs.auth_method,` +
		// oidc config
		` projections.apps4_oidc_configs.app_id,` +
		` projections.apps4_oidc_configs.version,` +
		` projections.apps4_oidc_configs.client_id,` +
		` projections.apps4_oidc_configs.redirect_uris,` +
		` projections.apps4_oidc_configs.response_types,` +
		` projections.apps4_oidc_configs.grant_types,` +
		` projections.apps4_oidc_configs.application_type,` +
		` projections.apps4_oidc_configs.auth_method_type,` +
		` projections.apps4_oidc_configs.post_logout_redirect_uris,` +
		` projections.apps4_oidc_configs.is_dev_mode,` +
		` projections.apps4_oidc_configs.access_token_type,` +
		` projections.apps4_oidc_configs.access_token_role_assertion,` +
		` projections.apps4_oidc_configs.id_token_role_assertion,` +
		` projections.apps4_oidc_configs.id_token_userinfo_assertion,` +
		` projections.apps4_oidc_configs.clock_skew,` +
		` projections.apps4_oidc_configs.additional_origins,` +
		//saml config
		` projections.apps4_saml_configs.app_id,` +
		` projections.apps4_saml_configs.entity_id,` +
		` projections.apps4_saml_configs.metadata,` +
		` projections.apps4_saml_configs.metadata_url,` +
		` COUNT(*) OVER ()` +
		` FROM projections.apps4` +
		` LEFT JOIN projections.apps4_api_configs ON projections.apps4.id = projections.apps4_api_configs.app_id AND projections.apps4.instance_id = projections.apps4_api_configs.instance_id` +
		` LEFT JOIN projections.apps4_oidc_configs ON projections.apps4.id = projections.apps4_oidc_configs.app_id AND projections.apps4.instance_id = projections.apps4_oidc_configs.instance_id` +
		` LEFT JOIN projections.apps4_saml_configs ON projections.apps4.id = projections.apps4_saml_configs.app_id AND projections.apps4.instance_id = projections.apps4_saml_configs.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedAppIDsQuery = regexp.QuoteMeta(`SELECT projections.apps4_api_configs.client_id,` +
		` projections.apps4_oidc_configs.client_id` +
		` FROM projections.apps4` +
		` LEFT JOIN projections.apps4_api_configs ON projections.apps4.id = projections.apps4_api_configs.app_id AND projections.apps4.instance_id = projections.apps4_api_configs.instance_id` +
		` LEFT JOIN projections.apps4_oidc_configs ON projections.apps4.id = projections.apps4_oidc_configs.app_id AND projections.apps4.instance_id = projections.apps4_oidc_configs.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedProjectIDByAppQuery = regexp.QuoteMeta(`SELECT projections.apps4.project_id` +
		` FROM projections.apps4` +
		` LEFT JOIN projections.apps4_api_configs ON projections.apps4.id = projections.apps4_api_configs.app_id AND projections.apps4.instance_id = projections.apps4_api_configs.instance_id` +
		` LEFT JOIN projections.apps4_oidc_configs ON projections.apps4.id = projections.apps4_oidc_configs.app_id AND projections.apps4.instance_id = projections.apps4_oidc_configs.instance_id` +
		` LEFT JOIN projections.apps4_saml_configs ON projections.apps4.id = projections.apps4_saml_configs.app_id AND projections.apps4.instance_id = projections.apps4_saml_configs.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	expectedProjectByAppQuery = regexp.QuoteMeta(`SELECT projections.projects3.id,` +
		` projections.projects3.creation_date,` +
		` projections.projects3.change_date,` +
		` projections.projects3.resource_owner,` +
		` projections.projects3.state,` +
		` projections.projects3.sequence,` +
		` projections.projects3.name,` +
		` projections.projects3.project_role_assertion,` +
		` projections.projects3.project_role_check,` +
		` projections.projects3.has_project_check,` +
		` projections.projects3.private_labeling_setting` +
		` FROM projections.projects3` +
		` JOIN projections.apps4 ON projections.projects3.id = projections.apps4.project_id AND projections.projects3.instance_id = projections.apps4.instance_id` +
		` LEFT JOIN projections.apps4_api_configs ON projections.apps4.id = projections.apps4_api_configs.app_id AND projections.apps4.instance_id = projections.apps4_api_configs.instance_id` +
		` LEFT JOIN projections.apps4_oidc_configs ON projections.apps4.id = projections.apps4_oidc_configs.app_id AND projections.apps4.instance_id = projections.apps4_oidc_configs.instance_id` +
		` LEFT JOIN projections.apps4_saml_configs ON projections.apps4.id = projections.apps4_saml_configs.app_id AND projections.apps4.instance_id = projections.apps4_saml_configs.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)

	appCols = database.StringArray{
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
		//saml config
		"app_id",
		"entity_id",
		"metadata",
		"metadata_url",
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
							// saml config
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
							// saml config
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
		}, {
			name:    "prepareAppsQuery saml app",
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
							// saml config
							"app-id",
							"https://test.com/saml/metadata",
							[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
							"https://test.com/saml/metadata",
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
						SAMLConfig: &SAMLApp{
							Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
							MetadataURL: "https://test.com/saml/metadata",
							EntityID:    "https://test.com/saml/metadata",
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           database.StringArray{"https://redirect.to/me"},
							ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
							IsDevMode:              true,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      database.StringArray{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							false,
							false,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           database.StringArray{"https://redirect.to/me"},
							ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
							IsDevMode:              false,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  false,
							AssertIDTokenRole:      false,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      database.StringArray{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							false,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           database.StringArray{"https://redirect.to/me"},
							ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
							IsDevMode:              true,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      false,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      database.StringArray{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							false,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           database.StringArray{"https://redirect.to/me"},
							ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
							IsDevMode:              false,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  false,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      database.StringArray{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
						OIDCConfig: &OIDCApp{
							Version:                domain.OIDCVersionV1,
							ClientID:               "oidc-client-id",
							RedirectURIs:           database.StringArray{"https://redirect.to/me"},
							ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
							IsDevMode:              false,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      database.StringArray{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
							nil,
							nil,
							nil,
							nil,
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
							// saml config
							nil,
							nil,
							nil,
							nil,
						},
						{
							"saml-app-id",
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
							// saml config
							"saml-app-id",
							"https://test.com/saml/metadata",
							[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
							"https://test.com/saml/metadata",
						},
					},
				),
			},
			object: &Apps{
				SearchResponse: SearchResponse{
					Count: 3,
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
							RedirectURIs:           database.StringArray{"https://redirect.to/me"},
							ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							AppType:                domain.OIDCApplicationTypeUserAgent,
							AuthMethodType:         domain.OIDCAuthMethodTypeNone,
							PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
							IsDevMode:              true,
							AccessTokenType:        domain.OIDCTokenTypeJWT,
							AssertAccessTokenRole:  true,
							AssertIDTokenRole:      true,
							AssertIDTokenUserinfo:  true,
							ClockSkew:              1 * time.Second,
							AdditionalOrigins:      database.StringArray{"additional.origin"},
							ComplianceProblems:     nil,
							AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
					{
						ID:            "saml-app-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						State:         domain.AppStateActive,
						Sequence:      20211109,
						Name:          "app-name",
						ProjectID:     "project-id",
						SAMLConfig: &SAMLApp{
							Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
							MetadataURL: "https://test.com/saml/metadata",
							EntityID:    "https://test.com/saml/metadata",
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
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
						// saml config
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
							// saml config
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           database.StringArray{"https://redirect.to/me"},
					ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      database.StringArray{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
				},
			},
		}, {
			name:    "prepareAppQuery saml app",
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
							// saml config
							"app-id",
							"https://test.com/saml/metadata",
							[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
							"https://test.com/saml/metadata",
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
				SAMLConfig: &SAMLApp{
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "https://test.com/saml/metadata",
					EntityID:    "https://test.com/saml/metadata",
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							false,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           database.StringArray{"https://redirect.to/me"},
					ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
					IsDevMode:              false,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      database.StringArray{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							false,
							true,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           database.StringArray{"https://redirect.to/me"},
					ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  false,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      database.StringArray{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							false,
							true,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           database.StringArray{"https://redirect.to/me"},
					ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      false,
					AssertIDTokenUserinfo:  true,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      database.StringArray{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
							database.StringArray{"https://redirect.to/me"},
							database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
							database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
							domain.OIDCApplicationTypeUserAgent,
							domain.OIDCAuthMethodTypeNone,
							database.StringArray{"post.logout.ch"},
							true,
							domain.OIDCTokenTypeJWT,
							true,
							true,
							false,
							1 * time.Second,
							database.StringArray{"additional.origin"},
							// saml config
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
				OIDCConfig: &OIDCApp{
					Version:                domain.OIDCVersionV1,
					ClientID:               "oidc-client-id",
					RedirectURIs:           database.StringArray{"https://redirect.to/me"},
					ResponseTypes:          database.EnumArray[domain.OIDCResponseType]{domain.OIDCResponseTypeIDTokenToken},
					GrantTypes:             database.EnumArray[domain.OIDCGrantType]{domain.OIDCGrantTypeImplicit},
					AppType:                domain.OIDCApplicationTypeUserAgent,
					AuthMethodType:         domain.OIDCAuthMethodTypeNone,
					PostLogoutRedirectURIs: database.StringArray{"post.logout.ch"},
					IsDevMode:              true,
					AccessTokenType:        domain.OIDCTokenTypeJWT,
					AssertAccessTokenRole:  true,
					AssertIDTokenRole:      true,
					AssertIDTokenUserinfo:  false,
					ClockSkew:              1 * time.Second,
					AdditionalOrigins:      database.StringArray{"additional.origin"},
					ComplianceProblems:     nil,
					AllowedOrigins:         database.StringArray{"https://redirect.to", "additional.origin"},
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
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
			name:    "prepareClientIDsQuery no result",
			prepare: prepareClientIDsQuery,
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
			name:    "prepareClientIDsQuery one result",
			prepare: prepareClientIDsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppIDsQuery,
					database.StringArray{"client_id", "client_id"},
					[][]driver.Value{
						{
							"app-id",
							nil,
						},
					},
				),
			},
			object: []string{"app-id"},
		},
		{
			name:    "prepareClientIDsQuery multiple result",
			prepare: prepareClientIDsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedAppIDsQuery,
					database.StringArray{"client_id", "client_id"},
					[][]driver.Value{
						{
							nil,
							"oidc-app-id",
						},
						{
							"api-app-id",
							nil,
						},
					},
				),
			},
			object: []string{"oidc-app-id", "api-app-id"},
		},
		{
			name:    "prepareClientIDsQuery sql err",
			prepare: prepareClientIDsQuery,
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
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
					database.StringArray{"project_id"},
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
